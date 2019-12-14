package collector

import (
	"log"

	"github.com/linode/linodego"
	"github.com/prometheus/client_golang/prometheus"
)

// KubernetesCollector represents a Linode Kubernetes Engine cluster (aka "LKE")
type KubernetesCollector struct {
	client linodego.Client
}

// NewKubernetesCollector creates a KubernetesCollector
func NewKubernetesCollector(client linodego.Client) *KubernetesCollector {
	// labelKeys := []string{"id", "label", "region"}
	return &KubernetesCollector{
		client: client,
	}
}

// Collect implements Collector interface and is called by Prometheus to collect metrics
func (c *KubernetesCollector) Collect(ch chan<- prometheus.Metric) {
	log.Println("[KubernetesCollector:Collect] Not yet implemented by linedogo")
}

// Describe implements Collector interface and is called by Prometheus to describe metrics
func (c *KubernetesCollector) Describe(ch chan<- *prometheus.Desc) {
	log.Println("[KubernetesCollector:Describe] Not yet implemented by linodego")
}
