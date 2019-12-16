# [Prometheus Exporter](https://prometheus.io/docs/instrumenting/exporters/) for [Linode](https://www.linode.com)

Initial implementation: does not do much!

Inspired by and templated from [DigitalOcean Exporter](https://github.com/metalmatze/digitalocean_exporter).

Thanks [metalmatze](https://github.com/metalmatze)!

## Development Installation

```bash
go get github.com/DazWilkin/linode-exporter
```
Then:
```bash
LINODE_TOKEN=[[YOUR-LINODE-API-TOKEN]]
ENDPOINT=":2112"
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
PORT=2112
docker run \
--interactive \
--tty \
--publish=${PORT}:2112 \
dazwilkin/linode-exporter:0d20258d4df27e592e8781a58640b363499b5b1c \
  --linode_token=${LINODE_TOKEN}
```
Or:
```bash
docker-compose --file=${PWD}/docker-compose.yaml up
```
## Metrics

| Name                                       | Type  | Description
| ----                                       | ----  | -----------
| linode_account_balance                     | Gauge ||
| linode_account_uninvoiced                  | Gauge ||
| linode_instance_count                      | Gauge ||
| linode_instance_cpu_average_utilization    | Gauge ||
| linode_instance_cpu_max_utilization        | Gauge ||
| linode_instance_io_total_blocks            | Gauge ||
| linode_instance_io_average_blocks          | Gauge ||
| linode_instance_swap_total_blocks          | Gauge ||
| linode_instance_swap_average_blocks        | Gauge ||
| linode_instance_ipv4_total_bits_received   | Gauge ||
| linode_instance_ipv4_average_bits_received | Gauge ||
| linode_instance_ipv4_total_bits_sent       | Gauge ||
| linode_instance_ipv4_average_bits_sent     | Gauge ||
| linode_instance_cpu_max_utilization        | Gauge ||
| linode_nodebalancer_count                  | Gauge ||
| linode_nodebalancer_transfer_total_bytes   | Gauge ||
| linode_nodebalancer_transfer_out_bytes     | Gauge ||
| linode_nodebalancer_transfer_in_bytes      | Gauge ||
| linode_support_tickets_count               | Gauge ||

Yeah, basic :-)

## Development

Each 'collector' is defined under `/collectors/[name].go`.

Collectors are instantiated by `main.go` with `registry.MustRegister(NewSomethingCollector(linodeClient))`

The `[name].go` collector implements Prometheus' Collector interface: `Collect` and `Describe`
