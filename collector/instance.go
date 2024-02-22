package collector

import (
	"context"
	"log"
	"strconv"
	"sync"

	"github.com/linode/linodego"
	"github.com/prometheus/client_golang/prometheus"
)

// InstanceCollector represents a Linode Instance (aka "Linode")
type InstanceCollector struct {
	client linodego.Client

	Up     *prometheus.Desc
	Disk   *prometheus.Desc
	Memory *prometheus.Desc
	CPUs   *prometheus.Desc

	//TODO(dazwilkin) IO swap
}

// NewInstanceCollector creates an InstanceCollector
func NewInstanceCollector(client linodego.Client) *InstanceCollector {
	log.Println("[NewInstanceCollector] Entered")
	subsystem := "instance"
	labelKeys := []string{"id", "label", "region"}
	return &InstanceCollector{
		client: client,

		Up: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "up"),
			"Status of Linode",
			labelKeys,
			nil,
		),
		Disk: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "disk"),
			"The amount of disk space in MB",
			labelKeys,
			nil,
		),
		Memory: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "cpu_max_utilization"),
			"The amount of RAM in MB",
			labelKeys,
			nil,
		),
		CPUs: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "io_total_blocks"),
			"The number of vCPUs",
			labelKeys,
			nil,
		),
	}
}

// Collect implements Collector interface and is called by Prometheus to collect metrics
func (c *InstanceCollector) Collect(ch chan<- prometheus.Metric) {
	log.Println("[InstanceCollector:Collect] Entered")
	ctx := context.Background()

	instances, err := c.client.ListInstances(ctx, nil)
	if err != nil {
		log.Println(err)
		return
	}
	log.Printf("[InstaneCollector:Collect] len(instances)=%d", len(instances))

	var wg sync.WaitGroup
	for _, instance := range instances {
		log.Printf("[InstanceCollector:Collect] Linode ID (%d)", instance.ID)

		wg.Add(1)
		go func(i linodego.Instance) {
			defer wg.Done()
			log.Printf("[InstanceCollector:Collect:go] Linode ID (%d)", i.ID)
			labelValues := []string{
				strconv.Itoa(i.ID),
				i.Label,
				i.Region,
			}

			ch <- prometheus.MustNewConstMetric(
				c.Up,
				prometheus.CounterValue,
				1.0,
				labelValues...,
			)

			ch <- prometheus.MustNewConstMetric(
				c.Disk,
				prometheus.GaugeValue,
				float64(i.Specs.Disk),
				labelValues...,
			)
			ch <- prometheus.MustNewConstMetric(
				c.Memory,
				prometheus.GaugeValue,
				float64(i.Specs.Memory),
				labelValues...,
			)

			ch <- prometheus.MustNewConstMetric(
				c.CPUs,
				prometheus.GaugeValue,
				float64(i.Specs.VCPUs),
				labelValues...,
			)

		}(instance)
	}
	wg.Wait()
	log.Println("[InstanceCollector:Collect] Completes")
}

// Describe implements Collector interface and is called by Prometheus to describe metrics
func (c *InstanceCollector) Describe(ch chan<- *prometheus.Desc) {
	log.Println("[InstanceCollector:Describe] Entered")
	ch <- c.Up
	ch <- c.Disk
	ch <- c.Memory
	ch <- c.CPUs
	log.Println("[InstanceCollector:Describe] Completes")
}
