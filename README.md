# carbon-tooling
Tooling to manage the Carbon stack. <br>
Injector component simulate the injection of metrics to carbon-relay-ng. <br>
Sink component simulate the receiving part. <br>
Prometheus observe the latency and metrics loss of the carbon-relay-ng thanks to the metrics of Injector and Sink.

## How to use it ?
Install docker <br>
Build a carbon-relay-ng image as criteo-carbon-relay-ng
```
docker build . -t criteo-carbon-relay-ng
```
Now you can run the the docker-compose to test the criteo-carbon-relay-ng
``` 
cd deployments/docker-compose/
docker-compose build
docker-compose up 
```
Check the prometheus alerts at http://localhost:9090/alerts
