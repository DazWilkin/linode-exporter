package collector

import (
	"context"
	"fmt"
	"log"
	"sync"

	"github.com/linode/linodego"
	"github.com/prometheus/client_golang/prometheus"
)

// ObjectStorageQuotaCollector represents a Linode object storage bucket
type ObjectStorageQuotaCollector struct {
	client linodego.Client

	BucketsUsage *prometheus.Desc
	BucketsLimit *prometheus.Desc
}

// NewObjectStorageQuotaCollector creates a ObjectStorageQuotaCollector
func NewObjectStorageQuotaCollector(client linodego.Client) *ObjectStorageQuotaCollector {
	log.Println("[NewObjectStorageQuotaCollector] Entered")
	subsystem := "objectstorage_quota"
	labelKeys := []string{"region", "endpoint_type", "endpoint"}
	return &ObjectStorageQuotaCollector{
		client: client,

		BucketsUsage: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "buckets_usage"),
			"Count of buckets in region",
			labelKeys,
			nil,
		),
		BucketsLimit: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "buckets_limit"),
			"Buckets limit in region",
			labelKeys,
			nil,
		),
	}
}

func (c *ObjectStorageQuotaCollector) collectBucketQuota(ctx context.Context, endpoint linodego.ObjectStorageEndpoint) ([]prometheus.Metric, error) {
	log.Printf("[ObjectStorageQuotaCollector:collectBucketQuota] Getting bucket usage and limit for (%s)", endpoint.Region)
	labelValues := []string{
		endpoint.Region,
		string(endpoint.EndpointType),
		*endpoint.S3Endpoint,
	}
	quotaID := fmt.Sprintf("obj-buckets-%s", *endpoint.S3Endpoint)
	log.Printf("[ObjectStorageQuotaCollector:collectBucketQuota] Getting bucket quota for (%s)", quotaID)
	bucketsQuota, err := c.client.GetObjectStorageQuotaUsage(ctx, quotaID)
	if err != nil {
		log.Printf("[ObjectStorageQuotaCollector:collectBucketQuota] Error getting bucket quota for (%s): %v", quotaID, err)
		return nil, err
	}
	return []prometheus.Metric{
		prometheus.MustNewConstMetric(
			c.BucketsUsage,
			prometheus.GaugeValue,
			float64(*bucketsQuota.Usage),
			labelValues...,
		),
		prometheus.MustNewConstMetric(
			c.BucketsLimit,
			prometheus.GaugeValue,
			float64(bucketsQuota.QuotaLimit),
			labelValues...,
		),
	}, nil
}

// Collect implements Collector interface and is called by Prometheus to collect metrics
func (c *ObjectStorageQuotaCollector) Collect(ch chan<- prometheus.Metric) {
	log.Println("[ObjectStorageQuotaCollector:Collect] Entered")
	ctx := context.Background()
	endpoints, err := GetLinodeObjectStorageEndpoints(c.client, ctx)
	log.Printf("[ObjectStorageQuotaCollector:Collect] len(endpoints)=%d", len(endpoints))
	if err != nil {
		log.Println(err)
		return
	}

	var wg sync.WaitGroup
	concurrencyLimit := 100
	sem := make(chan struct{}, concurrencyLimit)
	for _, endpoint := range endpoints {
		wg.Add(1)
		sem <- struct{}{} // acquire semaphore
		go func(endpoint linodego.ObjectStorageEndpoint) {
			defer wg.Done()
			defer func() { <-sem }() // release semaphore
			metrics, err := c.collectBucketQuota(ctx, endpoint)
			if err != nil {
				log.Printf("[ObjectStorageQuotaCollector:Collect] Error collecting bucket quota for (%s): %v", *endpoint.S3Endpoint, err)
				return
			}
			for _, metric := range metrics {
				ch <- metric
			}
		}(endpoint)
	}
	wg.Wait()
	log.Println("[ObjectStorageQuotaCollector:Collect] Completes")
}

// Describe implements Collector interface and is called by Prometheus to describe metrics
func (c *ObjectStorageQuotaCollector) Describe(ch chan<- *prometheus.Desc) {
	log.Println("[ObjectStorageQuotaCollector:Describe] Entered")
	ch <- c.BucketsUsage
	ch <- c.BucketsLimit
	log.Println("[ObjectStorageQuotaCollector:Describe] Completes")
}
