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

```bash
export LINODE_TOKEN=[[LINODE-API-TOKEN]]
```
Either:
```bash
PORT=9388
docker run \
--interactive \
--tty \
--publish=${PORT}:${PORT} \
dazwilkin/linode-exporter:0e4f56babde7b13aaec9abbb3585c9a0e188572d \
  --linode_token=${LINODE_TOKEN}
```
Or:
```bash
docker-compose --file=${PWD}/docker-compose.yaml up
```
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

## Port Allocation

Registered "Linode Exporter" on Prometheus Wiki's [Default Port Allocations](https://github.com/prometheus/prometheus/wiki/Default-port-allocations#exporters-starting-at-9100) with port 9388.