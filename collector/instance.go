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

	Count     *prometheus.Desc
	CPUAvg    *prometheus.Desc
	CPUMax    *prometheus.Desc
	IOSum     *prometheus.Desc
	IOAvg     *prometheus.Desc
	SwapSum   *prometheus.Desc
	SwapAvg   *prometheus.Desc
	IPv4RxSum *prometheus.Desc
	IPv4RxAvg *prometheus.Desc
	IPv4TxSum *prometheus.Desc
	IPv4TxAvg *prometheus.Desc

	//TODO(dazwilkin) IO swap
}

// NewInstanceCollector creates an InstanceCollector
func NewInstanceCollector(client linodego.Client) *InstanceCollector {
	log.Println("[NewInstanceCollector] Entered")
	fqName := name("instance")
	labelKeys := []string{"id", "label", "region"}
	return &InstanceCollector{
		client: client,

		Count: prometheus.NewDesc(
			fqName("count"),
			"Number of Linodes",
			labelKeys,
			nil,
		),
		CPUAvg: prometheus.NewDesc(
			fqName("cpu_average_utilization"),
			"CPU average utilization value for past 24 hours",
			labelKeys,
			nil,
		),
		CPUMax: prometheus.NewDesc(
			fqName("cpu_max_utilization"),
			"CPU max utilization value for past 24 hours",
			labelKeys,
			nil,
		),
		IOSum: prometheus.NewDesc(
			fqName("io_total_blocks"),
			"IO total blocks written in past 24 hours",
			labelKeys,
			nil,
		),
		IOAvg: prometheus.NewDesc(
			fqName("io_average_blocks"),
			"IO average blocks written in past 24 hours",
			labelKeys,
			nil,
		),
		SwapSum: prometheus.NewDesc(
			fqName("swap_total_blocks"),
			"Swap total blocks written in past 24 hours",
			labelKeys,
			nil,
		),
		SwapAvg: prometheus.NewDesc(
			fqName("swap_average_blocks"),
			"Swap average blocks written in past 24 hours",
			labelKeys,
			nil,
		),
		IPv4RxSum: prometheus.NewDesc(
			fqName("ipv4_total_bits_received"),
			"IPv4 total bits received over past 24 hours",
			labelKeys,
			nil,
		),
		IPv4RxAvg: prometheus.NewDesc(
			fqName("ipv4_average_bits_received"),
			"IPv4 average bits received over past 24 hours",
			labelKeys,
			nil,
		),
		IPv4TxSum: prometheus.NewDesc(
			fqName("ipv4_total_bits_sent"),
			"IPv4 total bits sent over past 24 hours",
			labelKeys,
			nil,
		),
		IPv4TxAvg: prometheus.NewDesc(
			fqName("ipv4_average_bits_sent"),
			"IPv4 average bits sent over past 24 hours",
			labelKeys,
			nil,
		),
	}
}

// Collect implements Collector interface and is called by Prometheus to collect metrics
func (c *InstanceCollector) Collect(ch chan<- prometheus.Metric) {
	log.Println("[InstanceCollector:Collect] Entered")
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	instances, err := c.client.ListInstances(ctx, nil)
	if err != nil {
		//TODO(dazwilkin) capture logs
		log.Println(err)
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
			log.Printf("[InstanceCollector:Collect:go] Linode ID (%d)", i.ID)
			labelValues := []string{
				fmt.Sprintf("%d", i.ID),
				i.Label,
				i.Region,
			}

			// https://developers.linode.com/api/v4/linode-instances-linode-id-stats
			// Appears (!) to be 64 values (always) [0] == epoch in ms? [1] == value
			var wg2 sync.WaitGroup
			log.Printf("[InstanceCollector:Collect:go] Linode ID (%d) -- getting stats", i.ID)
			stats, err := c.client.GetInstanceStats(ctx, i.ID)
			if err != nil {
				log.Println(err)
				return
			}

			// CPU
			wg2.Add(1)
			go func(i linodego.Instance) {
				defer wg2.Done()
				log.Printf("[InstanceCollector:Collect:go:go] Linode ID (%d) -- computing CPU stats", i.ID)
				ts := NewTimeSeries(stats.Data.CPU)

				ch <- prometheus.MustNewConstMetric(
					c.CPUAvg,
					prometheus.GaugeValue,
					ts.Avg(),
					labelValues...,
				)
				ch <- prometheus.MustNewConstMetric(
					c.CPUMax,
					prometheus.GaugeValue,
					ts.Max(),
					labelValues...,
				)
			}(i)

			// IO
			wg2.Add(1)
			go func(i linodego.Instance) {
				defer wg2.Done()
				log.Printf("[InstanceCollector:Collect:go:go] Linode ID (%d) -- computing IO stats", i.ID)
				ts := NewTimeSeries(stats.Data.IO.IO)

				ch <- prometheus.MustNewConstMetric(
					c.IOSum,
					prometheus.GaugeValue,
					ts.Sum(),
					labelValues...,
				)
				ch <- prometheus.MustNewConstMetric(
					c.IOAvg,
					prometheus.GaugeValue,
					ts.Avg(),
					labelValues...,
				)
			}(i)

			// Swap
			wg2.Add(1)
			go func(i linodego.Instance) {
				defer wg2.Done()
				log.Printf("[InstanceCollector:Collect:go:go] Linode ID (%d) -- computing Swap stats", i.ID)
				ts := NewTimeSeries(stats.Data.IO.Swap)

				ch <- prometheus.MustNewConstMetric(
					c.SwapSum,
					prometheus.GaugeValue,
					ts.Sum(),
					labelValues...,
				)
				ch <- prometheus.MustNewConstMetric(
					c.SwapAvg,
					prometheus.GaugeValue,
					ts.Avg(),
					labelValues...,
				)
			}(i)

			// IPv4 In
			wg2.Add(1)
			go func(i linodego.Instance) {
				defer wg2.Done()
				log.Printf("[InstanceCollector:Collect:go:go] Linode ID (%d) -- computing IPv4 Rx stats", i.ID)
				ts := NewTimeSeries(stats.Data.NetV4.In)

				ch <- prometheus.MustNewConstMetric(
					c.IPv4RxSum,
					prometheus.GaugeValue,
					ts.Sum(),
					labelValues...,
				)
				ch <- prometheus.MustNewConstMetric(
					c.IPv4RxAvg,
					prometheus.GaugeValue,
					ts.Avg(),
					labelValues...,
				)
			}(i)

			// IPv4 Out
			wg2.Add(1)
			go func(i linodego.Instance) {
				defer wg2.Done()
				log.Printf("[InstanceCollector:Collect:go:go] Linode ID (%d) -- computing IPv4 Tx stats", i.ID)
				ts := NewTimeSeries(stats.Data.NetV4.Out)

				ch <- prometheus.MustNewConstMetric(
					c.IPv4TxSum,
					prometheus.GaugeValue,
					ts.Sum(),
					labelValues...,
				)
				ch <- prometheus.MustNewConstMetric(
					c.IPv4TxAvg,
					prometheus.GaugeValue,
					ts.Avg(),
					labelValues...,
				)
			}(i)

			wg2.Wait()

		}(instance)
	}
	wg.Wait()
	log.Println("[InstanceCollector:Collect] Completes")
}

// Describe implements Collector interface and is called by Prometheus to describe metrics
func (c *InstanceCollector) Describe(ch chan<- *prometheus.Desc) {
	log.Println("[InstanceCollector:Describe] Entered")
	ch <- c.Count
	ch <- c.CPUAvg
	ch <- c.CPUMax
	ch <- c.IOSum
	ch <- c.IOAvg
	ch <- c.SwapSum
	ch <- c.SwapAvg
	ch <- c.IPv4RxSum
	ch <- c.IPv4RxAvg
	ch <- c.IPv4TxSum
	ch <- c.IPv4TxAvg
	log.Println("[InstanceCollector:Describe] Completes")
}
