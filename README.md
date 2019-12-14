# [Prometheus Exporter](https://prometheus.io/docs/instrumenting/exporters/) for [Linode](https://www.linode.com)

Initial implementation: does not do much!

Inspired by and templated from [DigitalOcean Exporter](https://github.com/metalmatze/digitalocean_exporter).

Thanks [metalmatze](https://github.com/metalmatze)!

## Installation

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

## Metrics

| Name                      | Type  | Description
| ----                      | ----  | -----------
| linode_instance_count     | Gauge ||
| linode_nodebalancer_count | Gauge ||

Yeah, basic :-)

