package collector

import (
	"context"
	"log"
	"strconv"
	"sync"

	"github.com/DazWilkin/linode-exporter/mock"
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
	fqName := name("kubernetes")
	return &KubernetesCollector{
		client: client,

		Up: prometheus.NewDesc(
			fqName("count"),
			"Status of Kubernetes cluster",
			[]string{"id", "label", "region", "version"},
			nil,
		),
		Pool: prometheus.NewDesc(
			fqName("pool"),
			"Size of Kubernetes node pool",
			[]string{"cluster_id", "id", "type"},
			nil,
		),
		Linode: prometheus.NewDesc(
			fqName("linode_up"),
			"Status of Kubernetes node pool Linode",
			[]string{"cluster_id", "pool_id", "id", "status"},
			nil,
		),
	}
}

// Collect implements Collector interface and is called by Prometheus to collect metrics
func (c *KubernetesCollector) Collect(ch chan<- prometheus.Metric) {
	log.Println("[KubernetesCollector:Collect] Entered")
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	// clusters, err := c.client.ListLKEClusters(ctx, nil)
	clusters, err := mock.NewClient().ListLKEClusters(ctx, nil)
	if err != nil {
		//TODO(dazwilkin) capture logs
		log.Println(err)
	}
	log.Printf("[KubernetesCollector:Collect] len(clusters)=%d", len(clusters))

	var wg sync.WaitGroup
	for _, cluster := range clusters {
		wg.Add(1)
		go func(k mock.LKECluster) {
			// go func(k linodego.LKECluster) {
			defer wg.Done()
			ch <- prometheus.MustNewConstMetric(
				c.Up,
				prometheus.CounterValue,
				1.0,
				// Label Values
				string(k.ID), k.Label, k.Region, k.Version,
			)
			// pools, err := c.client.ListLKEClusterPools(ctx, k.ID, nil)
			pools, err := mock.NewClient().ListLKEClusterPools(ctx, k.ID, nil)
			if err != nil {
				log.Println(err)
			}
			log.Printf("[KubernetesCollector:Collect] Cluster: %d len(nodepools)=%d", k.ID, len(pools))

			for _, pool := range pools {
				wg.Add(1)
				go func(p mock.LKEClusterPool) {
					// go func(p linodego.LKEClusterPool) {
					defer wg.Done()
					ch <- prometheus.MustNewConstMetric(
						c.Pool,
						prometheus.GaugeValue,
						float64(p.Count),
						// Label Values
						string(k.ID), string(p.ID), p.Type,
					)
					log.Printf("[KubernetesCollector:Collect] Cluster:%d Pool:%d", k.ID, p.ID)
					for _, l := range p.Linodes {
						ch <- prometheus.MustNewConstMetric(
							c.Linode,
							prometheus.CounterValue,
							// Metric will be 1 if LKELinodeReady, 0 otherwise
							// func(status linodego.LKELinodeStatus) (value float64) {
							func(status mock.LKELinodeStatus) (value float64) {
								// if status == linodego.LKELinodeReady {
								if status == mock.LKELinodeReady {
									log.Printf("[KubernetesCollector:Collect] Cluster:%d Pool:%d Linode:%d (%s)", k.ID, p.ID, l.ID, l.Status)
									value = 1.0
								}
								return value
							}(l.Status),
							// Label Values includes status string
							strconv.Itoa(k.ID), strconv.Itoa(p.ID), strconv.Itoa(*l.ID), string(l.Status),
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
