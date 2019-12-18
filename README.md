# [Prometheus Exporter](https://prometheus.io/docs/instrumenting/exporters/) for [Linode](https://www.linode.com)

Inspired by and templated from [DigitalOcean Exporter](https://github.com/metalmatze/digitalocean_exporter).

Thanks [metalmatze](https://github.com/metalmatze)!

![](images/linode_instance_count.png)

## Development Installation

```bash
go get github.com/DazWilkin/linode-exporter
```
Then:
```bash
LINODE_TOKEN=[[YOUR-LINODE-API-TOKEN]]
ENDPOINT=":9388"
PATH="/metrics"

go run github.com/DazWilkin/linode-exporter \
--linode_token=${LINODE_TOKEN} \
--endpoint=${ENDPOINT} \
--path=${PATH}
```

## Run-only Installation

### Linode Exporter only

Either:
```bash
LINODE_TOKEN=[[LINODE-API-TOKEN]]
PORT=9388
docker run \
--interactive \
--tty \
--publish=${PORT}:${PORT} \
dazwilkin/linode-exporter:3a761d146493f45b13bdb260a0238e6f5b98330b \
  --linode_token=${LINODE_TOKEN}
```

The exporter's metrics endpoint will be available on `http://localhost:${PORT}/metrics`

### Linode Exporter with [Prometheus](https://prometheus.io), [AlertManager](https://prometheus.io/docs/alerting/alertmanager/) and [cAdvisor](https://github.com/google/cadvisor)

**NB** AlertManager integration is a work-in-progress

The following 
```bash
LINODE_TOKEN=[[LINODE-API-TOKEN]]
docker-compose --file=${PWD}/docker-compose.yaml up
```
You may check the state of the services using:
```bash
docker-compose ps
```
And logs for a specific service using, e.g.:
```bash
docker-compose logs linode-exporter
```
The following endpoints are exposed:
+ Linode-Exporter metrics: `http://localhost:9388/metrics`
+ Prometheus UI: `http://localhost:9090`
+ AlertManager UI: `http://localhost:9093`
+ cAdvisor UI: `http://localhost:8085` 

**NB** cAdvisor is mapped to `:8085` rather than it's default port `:8080`

Using the Prometheus UI, you may begin querying metrics by typing `linode_` to see the available set.

The full list is below.

## Metrics

| Name                                       | Type  | Description
| ----                                       | ----  | -----------
| `linode_account_balance`                     | Gauge ||
| `linode_account_uninvoiced`                  | Gauge ||
| `linode_exporter_up`                         | Counter | A metric with a constant value of '1' labeled with go, OS and the exporter versions |
| `linode_instance_count`                      | Gauge ||
| `linode_instance_cpu_average_utilization`    | Gauge ||
| `linode_instance_cpu_max_utilization`        | Gauge ||
| `linode_instance_io_total_blocks`            | Gauge ||
| `linode_instance_io_average_blocks`          | Gauge ||
| `linode_instance_swap_total_blocks`          | Gauge ||
| `linode_instance_swap_average_blocks`        | Gauge ||
| `linode_instance_ipv4_total_bits_received`   | Gauge ||
| `linode_instance_ipv4_average_bits_received` | Gauge ||
| `linode_instance_ipv4_total_bits_sent`       | Gauge ||
| `linode_instance_ipv4_average_bits_sent`     | Gauge ||
| `linode_instance_cpu_max_utilization`        | Gauge ||
| `linode_nodebalancer_count`                  | Gauge ||
| `linode_nodebalancer_transfer_total_bytes`   | Gauge ||
| `linode_nodebalancer_transfer_out_bytes`     | Gauge ||
| `linode_nodebalancer_transfer_in_bytes`      | Gauge ||
| `linode_tickets_count`                       | Gauge ||

Please file issues and feature requests

## Development

Each 'collector' is defined under `/collectors/[name].go`.

Collectors are instantiated by `main.go` with `registry.MustRegister(NewSomethingCollector(linodeClient))`

The `[name].go` collector implements Prometheus' Collector interface: `Collect` and `Describe`

## Documentation

https://godoc.org/github.com/DazWilkin/linode-exporter/collector

## Port Allocation

Registered "Linode Exporter" on Prometheus Wiki's [Default Port Allocations](https://github.com/prometheus/prometheus/wiki/Default-port-allocations#exporters-starting-at-9100) with port 9388.