package collector

import (
	"context"
	"fmt"
	"log"

	"github.com/linode/linodego"
	"github.com/prometheus/client_golang/prometheus"
)

// InstanceCollector represents a Linode Instance (aka "Linode")
type InstanceCollector struct {
	client linodego.Client

	Total *prometheus.Desc
	CPU   *prometheus.Desc
}

// NewInstanceCollector creates an InstanceCollector
func NewInstanceCollector(client linodego.Client) *InstanceCollector {
	labels := []string{"id", "label", "region"}
	return &InstanceCollector{
		client: client,

		Total: prometheus.NewDesc(
			"linode_instance_count",
			"The total number of Linodes",
			labels,
			nil,
		),
		CPU: prometheus.NewDesc(
			"linode_instance_cpu_utilization",
			"The most recent CPU utilization value",
			labels,
			nil,
		),
	}
}
// Collect implements Collector interface and is called by Prometheus to collect metrics
func (c *InstanceCollector) Collect(ch chan<- prometheus.Metric) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	instances, err := c.client.ListInstances(ctx, nil)
	if err != nil {
		//TODO(dazwilkin) capture logs
		log.Fatal(err)
	}
	log.Printf("[main] len(instances)=%d", len(instances))

	ch <- prometheus.MustNewConstMetric(
		c.Total,
		prometheus.GaugeValue,
		float64(len(instances)),
		//TODO(dazwilkin) What metrics labels to use for this type of aggregate?
		[]string{"", "", ""}...,
	)

	for _, instance := range instances {
		labels := []string{
			fmt.Sprintf("%d", instance.ID),
			instance.Label,
			instance.Region,
		}
		stats, err := c.client.GetInstanceStats(ctx, instance.ID)
		if err != nil {
			log.Fatal(err)
		}

		// https://developers.linode.com/api/v4/linode-instances-linode-id-stats
		// for _, cpu := range stats.Data.CPU {
		ch <- prometheus.MustNewConstMetric(
			c.CPU,
			prometheus.GaugeValue,
			stats.Data.CPU[0][0],
			labels...,
		)
		// }

	}

}
// Describe implements Collector interface and is called by Prometheus to describe metrics
func (c *InstanceCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.Total
	ch <- c.CPU
}
