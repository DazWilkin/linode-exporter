package collector

import (
	"context"
	"fmt"
	"log"
	"sync"

	"github.com/linode/linodego"
	"github.com/prometheus/client_golang/prometheus"
)

// ObjectStorageQuotaCollector represents a Linode object storage bucket quota collector.
// It implements the prometheus.Collector interface.
type ObjectStorageQuotaCollector struct {
	client linodego.Client

	BucketsUsage *prometheus.Desc
	BucketsLimit *prometheus.Desc
}

const objectStorageConcurrencyLimit = 5

// NewObjectStorageQuotaCollector creates a new ObjectStorageQuotaCollector.
func NewObjectStorageQuotaCollector(client linodego.Client) *ObjectStorageQuotaCollector {
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

// buildLabelValues constructs label values for a given endpoint.
func buildLabelValues(endpoint linodego.ObjectStorageEndpoint) []string {
	return []string{
		endpoint.Region,
		string(endpoint.EndpointType),
		*endpoint.S3Endpoint,
	}
}

// collectBucketQuota collects bucket quota metrics for a single endpoint.
func (c *ObjectStorageQuotaCollector) collectBucketQuota(ctx context.Context, endpoint linodego.ObjectStorageEndpoint) ([]prometheus.Metric, error) {
	labelValues := buildLabelValues(endpoint)
	quotaID := fmt.Sprintf("obj-buckets-%s", *endpoint.S3Endpoint)
	bucketsQuota, err := c.client.GetObjectStorageQuotaUsage(ctx, quotaID)
	if err != nil {
		return nil, fmt.Errorf("error getting bucket quota for %s: %w", quotaID, err)
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

// Collect implements the prometheus.Collector interface and is called by Prometheus to collect metrics.
func (c *ObjectStorageQuotaCollector) Collect(ch chan<- prometheus.Metric) {
	ctx := context.Background()
	endpoints, err := GetLinodeObjectStorageEndpoints(c.client, ctx)
	if err != nil {
		log.Printf("[ObjectStorageQuotaCollector:Collect] failed to get endpoints: %v", err)
		return
	}

	var wg sync.WaitGroup
	sem := make(chan struct{}, objectStorageConcurrencyLimit)
	for _, endpoint := range endpoints {
		wg.Add(1)
		sem <- struct{}{} // acquire semaphore
		go func(endpoint linodego.ObjectStorageEndpoint) {
			defer wg.Done()
			defer func() { <-sem }() // release semaphore
			metrics, err := c.collectBucketQuota(ctx, endpoint)
			if err != nil {
				log.Printf("[ObjectStorageQuotaCollector:Collect] %v", err)
				return
			}
			for _, metric := range metrics {
				ch <- metric
			}
		}(endpoint)
	}
	wg.Wait()
}

// Describe implements the prometheus.Collector interface and is called by Prometheus to describe metrics.
func (c *ObjectStorageQuotaCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.BucketsUsage
	ch <- c.BucketsLimit
}
