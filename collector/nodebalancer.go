package collector

import (
	"context"
	"log"

	"github.com/linode/linodego"
	"github.com/prometheus/client_golang/prometheus"
)

// NodeBalancerCollector represents a Linode NodeBalancer
type NodeBalancerCollector struct {
	client linodego.Client

	Total *prometheus.Desc
}

// NewNodeBalancerCollector creates a NodeBalancerCollector
func NewNodeBalancerCollector(client linodego.Client) *NodeBalancerCollector {
	labels := []string{"id", "label", "region"}
	return &NodeBalancerCollector{
		client: client,

		Total: prometheus.NewDesc(
			"linode_nodebalancer_count",
			"The total number of NodeBalancers",
			labels,
			nil,
		),
	}
}

// Collect implements Collector interface and is called by Prometheus to collect metrics
func (c *NodeBalancerCollector) Collect(ch chan<- prometheus.Metric) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	nodebalancers, err := c.client.ListNodeBalancers(ctx, nil)
	if err != nil {
		//TODO(dazwilkin) capture logs: Loki?
		log.Fatal(err)
	}
	log.Printf("[main] len(nodebalancers)=%d", len(nodebalancers))

	ch <- prometheus.MustNewConstMetric(
		c.Total,
		prometheus.GaugeValue,
		float64(len(nodebalancers)),
		//TODO(dazwilkin) What metrics labels to use for this type of aggregate?
		[]string{"", "", ""}...,
	)
	// for _, nodebalancer := range nodebalancers {
	// 	labels := []string{
	// 		fmt.Sprintf("%s", nodebalancer.ID),
	// 	}
	// }
}

// Describe implements Collector interface and is called by Prometheus to describe metrics
func (c *NodeBalancerCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.Total
}
