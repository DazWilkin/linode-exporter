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
TOKEN=[[YOUR-LINODE-API-TOKEN]]
ENDPOINT=":2112"
PATH="/metrics"

go run github.com/DazWilkin/linode-exporter \
--token=${TOKEN} \
--endpoint=${ENDPOINT} \
--path=${PATH}
```

## Run-only Installation

```bash
export TOKEN=[[LINODE-API-TOKEN]]
docker-compose --file=${PWD}/docker-compose.yaml up
```
## Metrics

| Name                                       | Type  | Description
| ----                                       | ----  | -----------
| linode_instance_count                      | Gauge ||
| linode_instance_cpu_average_utilization    | Gauge ||
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

Yeah, basic :-)

