package collector

import (
	"context"
	"fmt"
	"log"

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
	labelKeys := []string{"id", "label", "region"}
	return &NodeBalancerCollector{
		client: client,

		Count: prometheus.NewDesc(
			"linode_nodebalancer_count",
			"The total number of NodeBalancers",
			labelKeys,
			nil,
		),
		TransferTotal: prometheus.NewDesc(
			"linode_nodebalancer_transfer_total_bytes",
			"The total number of MB transferred this month by the NodeBalancer",
			labelKeys,
			nil,
		),
		TransferOut: prometheus.NewDesc(
			"linode_nodebalancer_transfer_out_bytes",
			"The total number of MB transferred out this month by the NodeBalancer",
			labelKeys,
			nil,
		),
		TransferIn: prometheus.NewDesc(
			"linode_nodebalancer_transfer_in_bytes",
			"The total number of MB transferred in this month by the NodeBalancer",
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
		log.Fatal(err)
	}
	log.Printf("[NodeBalancerCollector:Collect] len(nodebalancers)=%d", len(nodebalancers))

	ch <- prometheus.MustNewConstMetric(
		c.Count,
		prometheus.GaugeValue,
		float64(len(nodebalancers)),
		//TODO(dazwilkin) What metrics labels to use for this type of aggregate?
		[]string{"", "", ""}...,
	)
	for _, nodebalancer := range nodebalancers {
		log.Printf("[NodeBalancerCollector:Collect] NodeBalancer ID (%d)", nodebalancer.ID)
		labelValues := []string{
			fmt.Sprintf("%d", nodebalancer.ID),
			*nodebalancer.Label,
			//TODO(dazwilkin) NodeBalanacer includes Tags too but these appear not key=value pairs
			nodebalancer.Region,
		}
		//TODO(dazwilkin) GetNodeBalancerStats is not implemented
		// stats, err := c.client.GetNodeBalancerStats(ctx, nodebalancer.ID)
		ch <- prometheus.MustNewConstMetric(
			c.TransferTotal,
			prometheus.GaugeValue,
			*nodebalancer.Transfer.Total,
			labelValues...,
		)
		ch <- prometheus.MustNewConstMetric(
			c.TransferOut,
			prometheus.GaugeValue,
			*nodebalancer.Transfer.Out,
			labelValues...,
		)
		ch <- prometheus.MustNewConstMetric(
			c.TransferIn,
			prometheus.GaugeValue,
			*nodebalancer.Transfer.In,
			labelValues...,
		)

	}
}

// Describe implements Collector interface and is called by Prometheus to describe metrics
func (c *NodeBalancerCollector) Describe(ch chan<- *prometheus.Desc) {
	log.Println("[NodeBalancerCollector:Describe] Entered")
	ch <- c.Count
	ch <- c.TransferTotal
	ch <- c.TransferOut
	ch <- c.TransferIn
}
