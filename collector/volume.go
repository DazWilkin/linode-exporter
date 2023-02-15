package collector

import (
	"context"
	"log"
	"strconv"
	"sync"

	"github.com/linode/linodego"
	"github.com/prometheus/client_golang/prometheus"
)

// VolumeCollector represents a Linode Volume
type VolumeCollector struct {
	client linodego.Client

	Up *prometheus.Desc
}

// NewVolumeCollector creates a new VolumeCollector
func NewVolumeCollector(client linodego.Client) *VolumeCollector {
	log.Println("[VolumeCollector] Entered")
	subsystem := "volume"
	return &VolumeCollector{
		client: client,

		Up: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "up"),
			"Status of Volume",
			[]string{"id", "label", "status", "region"},
			nil,
		),
	}
}

// Collect implements Collector interface and is called by Prometheus to collect metrics
func (c *VolumeCollector) Collect(ch chan<- prometheus.Metric) {
	log.Println("[VolumeCollector:Collect] Entered")
	ctx := context.Background()

	volumes, err := c.client.ListVolumes(ctx, nil)
	if err != nil {
		log.Println(err)
	}
	log.Printf("[VolumeCollector:Collect] len(volumes)=%d", len(volumes))

	var wg sync.WaitGroup
	for _, volume := range volumes {
		wg.Add(1)
		go func(v linodego.Volume) {
			defer wg.Done()
			ch <- prometheus.MustNewConstMetric(
				c.Up,
				prometheus.CounterValue,
				1.0,
				//Label Values
				strconv.Itoa(v.ID), v.Label, string(v.Status), v.Region,
			)
		}(volume)
	}
	wg.Wait()
	log.Println("[VolumeCollector:Collect] Completes")
}

// Describe implement Collector interface and is called by Prometheus to describe metrics
func (c *VolumeCollector) Describe(ch chan<- *prometheus.Desc) {
	log.Println("[VolumeCollector:Describe] Entered")
	ch <- c.Up
	log.Println("[VolumeCollector:Describe] Completes")
}
