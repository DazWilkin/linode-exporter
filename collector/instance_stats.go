package collector

import (
	"context"
	"log"
	"strconv"
	"sync"

	"github.com/linode/linodego"
	"github.com/prometheus/client_golang/prometheus"
)

// InstanceStatsCollector represents a Linode Instance (aka "Linode") Stats
type InstanceStatsCollector struct {
	client linodego.Client

	CPUUsage   *prometheus.Desc
	DiskIO     *prometheus.Desc
	NetworkIn  *prometheus.Desc
	NetworkOut *prometheus.Desc
}

// NewInstanceStatsCollector creates an InstanceStatsCollector
func NewInstanceStatsCollector(client linodego.Client) *InstanceStatsCollector {
	log.Println("[NewInstanceStatsCollector] Entered")
	subsystem := "instance_stats"

	return &InstanceStatsCollector{
		client: client,

		CPUUsage: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "cpu_usage"),
			"CPU usage percentage for Linode",
			[]string{"linode_id", "label", "region"},
			nil,
		),
		DiskIO: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "diskio"),
			"Disk IO operations for Linode",
			[]string{"linode_id", "label", "region", "type"},
			nil,
		),
		NetworkIn: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "network_in"),
			"Network incoming bytes for Linode",
			[]string{"linode_id", "label", "region"},
			nil,
		),
		NetworkOut: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "network_out"),
			"Network outgoing bytes for Linode",
			[]string{"linode_id", "label", "region"},
			nil,
		),
	}
}

// Collect implements Collector interface and is called by Prometheus to collect metrics
func (c *InstanceStatsCollector) Collect(ch chan<- prometheus.Metric) {
	log.Println("[InstanceStatsCollector:Collect] Entered")
	ctx := context.Background()

	instances, err := c.client.ListInstances(ctx, nil)
	if err != nil {
		log.Println(err)
	}
	log.Printf("[InstanceStatsCollector:Collect] len(instances)=%d", len(instances))

	var wg sync.WaitGroup
	for _, instance := range instances {
		log.Printf("[InstanceStatsCollector:Collect] Linode ID (%d)", instance.ID)

		wg.Add(1)
		go func(i linodego.Instance) {
			defer wg.Done()
			log.Printf("[InstanceStatsCollector:Collect:go] Linode ID (%d)", i.ID)
			instanceID := strconv.Itoa(i.ID)
			labelValues := []string{
				instanceID,
				i.Label,
				i.Region,
			}

			is, err := c.client.GetInstanceStats(ctx, i.ID)
			if err != nil {
				log.Println(err)
				return
			}

			cpuUsage := is.Data.CPU[0][1]
			ch <- prometheus.MustNewConstMetric(
				c.CPUUsage,
				prometheus.GaugeValue,
				cpuUsage,
				labelValues...,
			)
			diskIO := is.Data.IO.IO[0][1]
			ch <- prometheus.MustNewConstMetric(
				c.DiskIO,
				prometheus.GaugeValue,
				diskIO,
				append(labelValues, "io")...,
			)
			swapIO := is.Data.IO.Swap[0][1]
			ch <- prometheus.MustNewConstMetric(
				c.DiskIO,
				prometheus.GaugeValue,
				swapIO,
				append(labelValues, "swap")...,
			)
			networkIn := is.Data.NetV4.In[0][1]
			ch <- prometheus.MustNewConstMetric(
				c.NetworkIn,
				prometheus.GaugeValue,
				networkIn,
				labelValues...,
			)
			networkOut := is.Data.NetV4.Out[0][1]
			ch <- prometheus.MustNewConstMetric(
				c.NetworkOut,
				prometheus.GaugeValue,
				networkOut,
				labelValues...,
			)

		}(instance)
	}
	wg.Wait()
	log.Println("[InstanceStatsCollector:Collect] Completes")
}

// Describe implements Collector interface and is called by Prometheus to describe metrics
func (c *InstanceStatsCollector) Describe(ch chan<- *prometheus.Desc) {
	log.Println("[InstanceStatsCollector:Describe] Entered")
	ch <- c.CPUUsage
	ch <- c.DiskIO
	ch <- c.NetworkIn
	ch <- c.NetworkOut
	log.Println("[InstanceStatsCollector:Describe] Completes")
}
