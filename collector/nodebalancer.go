package collector

import (
	"context"
	"fmt"
	"log"
	"sync"

	"github.com/linode/linodego"
	"github.com/prometheus/client_golang/prometheus"
)

// NodeBalancerCollector represents a Linode NodeBalancer
type NodeBalancerCollector struct {
	client linodego.Client

	Count         *prometheus.Desc
	TransferTotal *prometheus.Desc
	TransferOut   *prometheus.Desc
	TransferIn    *prometheus.Desc
}

// NewNodeBalancerCollector creates a NodeBalancerCollector
func NewNodeBalancerCollector(client linodego.Client) *NodeBalancerCollector {
	log.Println("[NewNodeBalancerCollector] Entered")
	subsystem := "nodebalancer"
	labelKeys := []string{"id", "label", "region"}
	return &NodeBalancerCollector{
		client: client,

		Count: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "count"),
			"Number of NodeBalancers",
			labelKeys,
			nil,
		),
		TransferTotal: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "transfer_total_bytes"),
			"MB transferred this month by the NodeBalancer",
			labelKeys,
			nil,
		),
		TransferOut: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "transfer_out_bytes"),
			"MB transferred out this month by the NodeBalancer",
			labelKeys,
			nil,
		),
		TransferIn: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "transfer_in_bytes"),
			"MB transferred in this month by the NodeBalancer",
			labelKeys,
			nil,
		),
	}
}

// Collect implements Collector interface and is called by Prometheus to collect metrics
func (c *NodeBalancerCollector) Collect(ch chan<- prometheus.Metric) {
	log.Println("[NodeBalancerCollector:Collect] Entered")
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	nodebalancers, err := c.client.ListNodeBalancers(ctx, nil)
	if err != nil {
		//TODO(dazwilkin) capture logs: Loki?
		log.Println(err)
	}
	log.Printf("[NodeBalancerCollector:Collect] len(nodebalancers)=%d", len(nodebalancers))

	ch <- prometheus.MustNewConstMetric(
		c.Count,
		prometheus.GaugeValue,
		float64(len(nodebalancers)),
		//TODO(dazwilkin) What metrics labels to use for this type of aggregate?
		[]string{"", "", ""}...,
	)

	var wg sync.WaitGroup
	for _, nodebalancer := range nodebalancers {
		log.Printf("[NodeBalancerCollector:Collect] NodeBalancer ID (%d)", nodebalancer.ID)
		wg.Add(1)
		go func(nb linodego.NodeBalancer) {
			defer wg.Done()
			log.Printf("[NodeBalancerCollector:Collect:go] NodeBalancer ID (%d)", nb.ID)
			labelValues := []string{
				fmt.Sprintf("%d", nb.ID),
				*nb.Label,
				//TODO(dazwilkin) NodeBalancer includes Tags too but these appear not key=value pairs
				nb.Region,
			}
			// nb.Transfer.[Total|Out|In] may be nil; only report these values when non-nil
			if nb.Transfer.Total != nil {
				ch <- prometheus.MustNewConstMetric(
					c.TransferTotal,
					prometheus.GaugeValue,
					*nb.Transfer.Total,
					labelValues...,
				)
			}
			if nb.Transfer.Out != nil {
				ch <- prometheus.MustNewConstMetric(
					c.TransferOut,
					prometheus.GaugeValue,
					*nb.Transfer.Out,
					labelValues...,
				)
			}
			if nb.Transfer.In != nil {
				ch <- prometheus.MustNewConstMetric(
					c.TransferIn,
					prometheus.GaugeValue,
					*nb.Transfer.In,
					labelValues...,
				)
			}

		}(nodebalancer)
	}
	wg.Wait()
	log.Println("[NodeBalancerCollector:Collect] Completes")
}

// Describe implements Collector interface and is called by Prometheus to describe metrics
func (c *NodeBalancerCollector) Describe(ch chan<- *prometheus.Desc) {
	log.Println("[NodeBalancerCollector:Describe] Entered")
	ch <- c.Count
	ch <- c.TransferTotal
	ch <- c.TransferOut
	ch <- c.TransferIn
	log.Println("[NodeBalancerCollector:Describe] Completes")
}
