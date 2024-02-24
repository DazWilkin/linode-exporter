package collector

import (
	"context"
	"log"
	"strconv"
	"sync"

	"github.com/linode/linodego"
	"github.com/prometheus/client_golang/prometheus"
)

// KubernetesCollector represents a Linode Kubernetes Engine cluster (aka "LKE")
type KubernetesCollector struct {
	client linodego.Client

	Up     *prometheus.Desc
	Pool   *prometheus.Desc
	Linode *prometheus.Desc
}

// NewKubernetesCollector creates a KubernetesCollector
func NewKubernetesCollector(client linodego.Client) *KubernetesCollector {
	log.Println("[NewKubernetesCollector] Entered")
	subsystem := "kubernetes"
	return &KubernetesCollector{
		client: client,

		Up: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "up"),
			"Status of Kubernetes cluster",
			[]string{"id", "label", "region", "version"},
			nil,
		),
		Pool: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "pool"),
			"Size of Kubernetes node pool",
			[]string{"cluster_id", "id", "type"},
			nil,
		),
		Linode: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "linode_up"),
			"Status of Kubernetes node pool Linode",
			[]string{"cluster_id", "pool_id", "id", "status"},
			nil,
		),
	}
}

// Collect implements Collector interface and is called by Prometheus to collect metrics
func (c *KubernetesCollector) Collect(ch chan<- prometheus.Metric) {
	log.Println("[KubernetesCollector:Collect] Entered")
	ctx := context.Background()

	clusters, err := c.client.ListLKEClusters(ctx, nil)
	if err != nil {
		log.Println(err)
		return
	}
	log.Printf("[KubernetesCollector:Collect] len(clusters)=%d", len(clusters))

	var wg sync.WaitGroup
	for _, cluster := range clusters {
		wg.Add(1)
		go func(k linodego.LKECluster) {
			defer wg.Done()
			ch <- prometheus.MustNewConstMetric(
				c.Up,
				prometheus.CounterValue,
				1.0,
				// Label Values
				strconv.Itoa(k.ID), k.Label, k.Region, k.K8sVersion,
			)
			pools, err := c.client.ListLKEClusterPools(ctx, k.ID, nil)
			if err != nil {
				log.Println(err)
				return
			}
			log.Printf("[KubernetesCollector:Collect] Cluster: %d len(nodepools)=%d", k.ID, len(pools))

			for _, pool := range pools {
				wg.Add(1)
				go func(p linodego.LKEClusterPool) {
					defer wg.Done()
					ch <- prometheus.MustNewConstMetric(
						c.Pool,
						prometheus.GaugeValue,
						float64(p.Count),
						// Label Values
						strconv.Itoa(k.ID), strconv.Itoa(p.ID), p.Type,
					)
					log.Printf("[KubernetesCollector:Collect] Cluster:%d Pool:%d", k.ID, p.ID)
					for _, l := range p.Linodes {
						ch <- prometheus.MustNewConstMetric(
							c.Linode,
							prometheus.CounterValue,
							// Metric will be 1 if LKELinodeReady, 0 otherwise
							// func(status linodego.LKELinodeStatus) (value float64) {
							func(status linodego.LKELinodeStatus) (value float64) {
								// if status == linodego.LKELinodeReady {
								if status == linodego.LKELinodeReady {
									log.Printf("[KubernetesCollector:Collect] Cluster:%d Pool:%d Linode:%s (%s)", k.ID, p.ID, l.ID, string(l.Status))
									value = 1.0
								}
								return value
							}(l.Status),
							// Label Values includes status string
							strconv.Itoa(k.ID), strconv.Itoa(p.ID), l.ID, string(l.Status),
						)
					}
				}(pool)
			}
		}(cluster)
	}
	wg.Wait()
	log.Println("[KubernetesCollector:Collect] Completes")

}

// Describe implements Collector interface and is called by Prometheus to describe metrics
func (c *KubernetesCollector) Describe(ch chan<- *prometheus.Desc) {
	log.Println("[KubernetesCollector:Describe] Entered")
	ch <- c.Up
	ch <- c.Pool
	ch <- c.Linode
	log.Println("[KubernetesCollector:Describe] Completes")
}
