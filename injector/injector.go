package main

import (
	"flag"
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
	"net"
	"net/http"
	"time"
)

var (
	metricsSent = promauto.NewCounter(prometheus.CounterOpts{
		Name: "injector_metrics_sent_total",
		Help: "The total number of metrics sent",
	})
	sendPause = flag.Duration("sendPause", 1*time.Second, "sleep between two send of metrics in seconds")
	host      = flag.String("host", "carbon-relay-ng", "destination host")
	port      = flag.Int("port", 2003, "destination port")
	log       = zap.NewExample()
)

func initTcpConnection(host string, port int) net.Conn {
	hostPort := fmt.Sprintf("%s:%d", host, port)
	log.Info("Trying to connect to carbon relay...", zap.String("destination", hostPort))
	conn, err := net.Dial("tcp", hostPort)
	retryConnectionMax := 10
	for retryConnection := 1; retryConnection <= retryConnectionMax && err != nil; retryConnection++ {
		conn, err = net.Dial("tcp", hostPort)
		log.Error("Fail to connect.", zap.String("destination", hostPort), zap.Int("retry", retryConnection), zap.Error(err))
		time.Sleep(1 * time.Second)
	}
	log.Info("Connected to server", zap.String("destination", hostPort))
	return conn
}

func getMetricMessage() string {
	metricPath := "my.custom.metric"
	nowNanoSec := time.Now().UnixNano()
	metricMessage := fmt.Sprintf("%s %d %d\n", metricPath, nowNanoSec, nowNanoSec/1e9)
	return metricMessage
}

func loopingCallSendMetrics() {
	go func() {
		log.Info("configuration dump", zap.Duration("sendPause", *sendPause), zap.String("host", *host), zap.Int("port", *port))
		conn := initTcpConnection(*host, *port)
		for {
			message := getMetricMessage()
			// send to socket
			_, err := fmt.Fprintf(conn, message)
			if err != nil {
				log.Error("Fail to send message", zap.String("message", message))
			}
			metricsSent.Inc()
			time.Sleep(*sendPause)
		}
	}()
}

func exposePrometheusMetricsHttpEndpoint() {
	http.Handle("/metrics", promhttp.Handler())
	prometheusPort := ":2112"
	err := http.ListenAndServe(prometheusPort, nil)
	if err != nil {
		log.Error("Prometheus fail to bind port", zap.String("prometheusPort", prometheusPort))
	}
}
func main() {
	flag.Parse()
	loopingCallSendMetrics()
	exposePrometheusMetricsHttpEndpoint()
}
