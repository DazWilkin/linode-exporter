package collector

import (
	"context"
	"fmt"
	"log"
	"sync"

	"github.com/linode/linodego"
	"github.com/prometheus/client_golang/prometheus"
)

// InstanceCollector represents a Linode Instance (aka "Linode")
type InstanceCollector struct {
	client linodego.Client

	Count  *prometheus.Desc
	CPUAve *prometheus.Desc
	CPUMax *prometheus.Desc
}

// NewInstanceCollector creates an InstanceCollector
func NewInstanceCollector(client linodego.Client) *InstanceCollector {
	labelKeys := []string{"id", "label", "region"}
	return &InstanceCollector{
		client: client,

		Count: prometheus.NewDesc(
			"linode_instance_count",
			"The total number of Linodes",
			labelKeys,
			nil,
		),
		CPUAve: prometheus.NewDesc(
			"linode_instance_cpu_average_utilization",
			"The most recent CPU average utilization value",
			labelKeys,
			nil,
		),
		CPUMax: prometheus.NewDesc(
			"linode_instance_cpu_max_utilization",
			"The most recent CPU max utilization value",
			labelKeys,
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
	log.Printf("[InstaneCollector:Collect] len(instances)=%d", len(instances))

	ch <- prometheus.MustNewConstMetric(
		c.Count,
		prometheus.GaugeValue,
		float64(len(instances)),
		//TODO(dazwilkin) What metrics labels to use for this type of aggregate?
		[]string{"", "", ""}...,
	)

	var wg sync.WaitGroup
	for _, instance := range instances {
		log.Printf("[InstanceCollector:Collect] Linode ID (%d)", instance.ID)

		wg.Add(1)
		go func(i linodego.Instance) {
			defer wg.Done()
			labelValues := []string{
				fmt.Sprintf("%d", i.ID),
				i.Label,
				i.Region,
			}

			// https://developers.linode.com/api/v4/linode-instances-linode-id-stats
			// Appears (!) to be 64 values (always) [0] == epoch in ms? [1] == value
			stats, err := c.client.GetInstanceStats(ctx, i.ID)
			if err != nil {
				log.Fatal(err)
			}

			var total, max float64
			for _, cpu := range stats.Data.CPU {
				if cpu[1] != 0.0 {
					total += cpu[1]
					if cpu[1] > max {
						max = cpu[1]
					}
				}
			}

			ch <- prometheus.MustNewConstMetric(
				c.CPUAve,
				prometheus.GaugeValue,
				total/float64(len(stats.Data.CPU)),
				labelValues...,
			)
			ch <- prometheus.MustNewConstMetric(
				c.CPUMax,
				prometheus.GaugeValue,
				max,
				labelValues...,
			)
		}(instance)
	}
	wg.Wait()
}

// Describe implements Collector interface and is called by Prometheus to describe metrics
func (c *InstanceCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.Count
	ch <- c.CPUAve
	ch <- c.CPUMax
}
