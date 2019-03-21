package main

import (
	"bufio"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"
)

var (
	metricsInjectorReceived = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "sink_metrics_received_total",
		Help: "The total number of metrics received",
	}, []string{"input"})
	metricsLatency = promauto.NewHistogram(prometheus.HistogramOpts{
		Name:    "sink_metrics_duration_seconds",
		Help:    "Histogram latency metrics",
		Buckets: prometheus.ExponentialBuckets(0.001, 10, 4),
	})
	log = zap.NewExample()
)

func getLatency(epochTimeStampInNanoSec int64) int64 {
	currentEpochTimeInNanoSec := time.Now().UnixNano()
	latency := currentEpochTimeInNanoSec - epochTimeStampInNanoSec
	return latency
}

func getLatencyInSecond(latencyInNanoSec int64) float64 {
	return float64(latencyInNanoSec / 1e9)
}

func isMetricFromCarbonRelayNg(metricPath string) bool {
	return strings.HasPrefix(metricPath, "service_is_carbon-relay-ng")
}

func handleMetricMessage(message string) {
	metricMessage := strings.Split(message, " ")
	if len(metricMessage) == 3 {
		if isMetricFromCarbonRelayNg(metricMessage[0]) {
			//log.Debug("metric received from injector", zap.String("metricPath", metricMessage[0]), zap.String("value", metricMessage[1]), zap.String("timestamp", metricMessage[2]))

			metricsInjectorReceived.With(prometheus.Labels{"input": "carbon-relay"}).Inc()

		} else {
			epochTimeStampInNanoSec, _ := strconv.ParseInt(metricMessage[1], 10, 64)
			latency := getLatency(epochTimeStampInNanoSec)
			//log.Debug("metric received from injector", zap.Int64("latency", latency), zap.String("metricPath", metricMessage[0]), zap.Int64("value", epochTimeStampInNanoSec), zap.String("timestamp", metricMessage[2]))
			metricsLatency.Observe(getLatencyInSecond(latency))
			metricsInjectorReceived.With(prometheus.Labels{"input": "injector"}).Inc()
		}
	}
}

func loopingCallReceiveMetrics() {
	go func() {
		ln, _ := net.Listen("tcp", ":2003")
		log.Info("Server ready waiting connection...")
		for {
			conn, _ := ln.Accept()
			//fmt.Println("Connection accepted")
			buffer := bufio.NewReader(conn)
			for {
				message, _ := buffer.ReadString('\n')
				message = strings.TrimSuffix(message, "\n")
				//log.Debug("message dump", zap.String("message", message))
				handleMetricMessage(message)
			}
		}
	}()
}

func exposePrometheusMetricsHttpEndpoint() {
	http.Handle("/metrics", promhttp.Handler())
	prometheusPort := ":2112"
	err := http.ListenAndServe(prometheusPort, nil)
	if err != nil {
		log.Error("Prometheus fail to bind to port", zap.String("port", prometheusPort))
	}
}

func main() {
	loopingCallReceiveMetrics()
	exposePrometheusMetricsHttpEndpoint()
}
