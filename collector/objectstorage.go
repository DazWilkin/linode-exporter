package collector

import (
	"context"
	"log"
	"sync"

	"github.com/linode/linodego"
	"github.com/prometheus/client_golang/prometheus"
)

// ObjectStorageCollector represents a Linode object storage bucket
type ObjectStorageCollector struct {
	client linodego.Client

	Size         *prometheus.Desc
	ObjectsCount *prometheus.Desc
}

// NewObjectStorageCollector creates a ObjectStorageCollector
func NewObjectStorageCollector(client linodego.Client) *ObjectStorageCollector {
	log.Println("[NewObjectStorageCollector] Entered")
	subsystem := "objectstorage"
	labelKeys := []string{"label", "region"}
	return &ObjectStorageCollector{
		client: client,

		Size: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "size_bytes"),
			"Size of a bucket (in bytes)",
			labelKeys,
			nil,
		),
		ObjectsCount: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "objects_count"),
			"Count of objects in a bucket",
			labelKeys,
			nil,
		),
	}
}

// Collect implements Collector interface and is called by Prometheus to collect metrics
func (c *ObjectStorageCollector) Collect(ch chan<- prometheus.Metric) {
	log.Println("[ObjectStorageCollector:Collect] Entered")
	ctx := context.Background()

	buckets, err := c.client.ListObjectStorageBuckets(ctx, nil)
	if err != nil {
		log.Println(err)
		return
	}
	log.Printf("[ObjectStorageCollector:Collect] len(buckets)=%d", len(buckets))

	var wg sync.WaitGroup
	for _, bucket := range buckets {
		log.Printf("[ObjectStorageCollector:Collect] Bucket ID (%s)", bucket.Label)

		wg.Add(1)
		go func(bucket linodego.ObjectStorageBucket) {
			defer wg.Done()
			log.Printf("[ObjectStorageCollector:Collect:go] Bucket ID (%s)", bucket.Label)
			labelValues := []string{
				bucket.Label,
				bucket.Region,
			}

			ch <- prometheus.MustNewConstMetric(
				c.Size,
				prometheus.GaugeValue,
				float64(bucket.Size),
				labelValues...,
			)
			ch <- prometheus.MustNewConstMetric(
				c.ObjectsCount,
				prometheus.GaugeValue,
				float64(bucket.Objects),
				labelValues...,
			)
		}(bucket)
	}
	wg.Wait()
	log.Println("[ObjectStorageCollector:Collect] Completes")
}

// Describe implements Collector interface and is called by Prometheus to describe metrics
func (c *ObjectStorageCollector) Describe(ch chan<- *prometheus.Desc) {
	log.Println("[ObjectStorageCollector:Describe] Entered")
	ch <- c.Size
	log.Println("[ObjectStorageCollector:Describe] Completes")
}
