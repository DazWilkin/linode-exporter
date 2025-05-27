package collector

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/linode/linodego"
	"github.com/prometheus/client_golang/prometheus"
)

// ObjectStorageQuotaCollector represents a Linode object storage bucket quota collector.
// It implements the prometheus.Collector interface.
type ObjectStorageQuotaCollector struct {
	client linodego.Client

	quotaTypeDescriptors map[string]struct {
		usage *prometheus.Desc
		limit *prometheus.Desc
	}

	endpointsCache      []linodego.ObjectStorageEndpoint
	endpointsCacheTime  time.Time
	endpointsCacheMutex sync.Mutex
}

const (
	objectStorageConcurrencyLimit = 50
	endpointsCacheExpiry          = 30 * time.Second
)

// NewObjectStorageQuotaCollector creates a new ObjectStorageQuotaCollector.
func NewObjectStorageQuotaCollector(client linodego.Client) *ObjectStorageQuotaCollector {
	subsystem := "objectstorage_quota"
	labelKeys := []string{"region", "endpoint_type", "endpoint"}
	quotaTypeDescriptors := map[string]struct {
		usage *prometheus.Desc
		limit *prometheus.Desc
	}{
		"buckets": {
			usage: prometheus.NewDesc(
				prometheus.BuildFQName(namespace, subsystem, "buckets_usage"),
				"Count of buckets in region",
				labelKeys,
				nil,
			),
			limit: prometheus.NewDesc(
				prometheus.BuildFQName(namespace, subsystem, "buckets_limit"),
				"Buckets limit in region",
				labelKeys,
				nil,
			),
		},
		"bytes": {
			usage: prometheus.NewDesc(
				prometheus.BuildFQName(namespace, subsystem, "bytes_usage"),
				"Count of bytes in region",
				labelKeys,
				nil,
			),
			limit: prometheus.NewDesc(
				prometheus.BuildFQName(namespace, subsystem, "bytes_limit"),
				"Bytes limit in region",
				labelKeys,
				nil,
			),
		},
	}
	return &ObjectStorageQuotaCollector{
		client:               client,
		quotaTypeDescriptors: quotaTypeDescriptors,
	}
}

// Collect implements the prometheus.Collector interface and is called by Prometheus to collect metrics.
func (c *ObjectStorageQuotaCollector) Collect(ch chan<- prometheus.Metric) {
	ctx := context.Background()
	endpoints, err := c.getCachedEndpoints(ctx)
	if err != nil {
		log.Printf("[ObjectStorageQuotaCollector:Collect] failed to get endpoints: %v", err)
		return
	}

	var wg sync.WaitGroup
	sem := make(chan struct{}, objectStorageConcurrencyLimit)
	for _, endpoint := range endpoints {
		for quotaType, descriptors := range c.quotaTypeDescriptors {
			wg.Add(1)
			sem <- struct{}{} // acquire semaphore
			go func(endpoint linodego.ObjectStorageEndpoint, quotaType string, usageDesc, limitDesc *prometheus.Desc) {
				defer wg.Done()
				defer func() { <-sem }() // release semaphore
				metrics, err := c.collectObjectStorageQuota(ctx, buildLabelValues(endpoint), *endpoint.S3Endpoint, quotaType, prometheus.GaugeValue, *usageDesc, *limitDesc)
				if err != nil {
					log.Printf("[ObjectStorageQuotaCollector:Collect] %v", err)
					return
				}
				for _, metric := range metrics {
					ch <- metric
				}
			}(endpoint, quotaType, descriptors.usage, descriptors.limit)
		}
	}
	wg.Wait()
}

// Describe implements the prometheus.Collector interface and is called by Prometheus to describe metrics.
func (c *ObjectStorageQuotaCollector) Describe(ch chan<- *prometheus.Desc) {
	for _, descriptors := range c.quotaTypeDescriptors {
		ch <- descriptors.usage
		ch <- descriptors.limit
	}
}

// collect bucket quota metrics for a single endpoint.
func (c *ObjectStorageQuotaCollector) collectObjectStorageQuota(ctx context.Context, labelValues []string, s3Endpoint string, quotaType string, metricType prometheus.ValueType, usageMetricDesc prometheus.Desc, limitMetricDesc prometheus.Desc) ([]prometheus.Metric, error) {
	quotaID := fmt.Sprintf("obj-%s-%s", quotaType, s3Endpoint)
	ObjectStorageQuota, err := c.client.GetObjectStorageQuotaUsage(ctx, quotaID)
	if err != nil {
		return nil, fmt.Errorf("error getting bucket quota for %s: %w", quotaID, err)
	}
	return []prometheus.Metric{
		prometheus.MustNewConstMetric(
			&usageMetricDesc,
			metricType,
			float64(*ObjectStorageQuota.Usage),
			labelValues...,
		),
		prometheus.MustNewConstMetric(
			&limitMetricDesc,
			metricType,
			float64(ObjectStorageQuota.QuotaLimit),
			labelValues...,
		),
	}, nil
}

// returns cached endpoints if not expired, otherwise fetches and updates the cache.
func (c *ObjectStorageQuotaCollector) getCachedEndpoints(ctx context.Context) ([]linodego.ObjectStorageEndpoint, error) {
	c.endpointsCacheMutex.Lock()
	defer c.endpointsCacheMutex.Unlock()
	if time.Since(c.endpointsCacheTime) < endpointsCacheExpiry && c.endpointsCache != nil {
		log.Printf("[ObjectStorageQuotaCollector:getCachedEndpoints] using cached endpoints")
		return c.endpointsCache, nil
	}
	log.Printf("[ObjectStorageQuotaCollector:getCachedEndpoints] fetching endpoints from API")
	endpoints, err := GetLinodeObjectStorageEndpoints(c.client, ctx)
	if err != nil {
		return nil, err
	}
	c.endpointsCache = endpoints
	c.endpointsCacheTime = time.Now()
	return endpoints, nil
}

// constructs label values for a given endpoint.
func buildLabelValues(endpoint linodego.ObjectStorageEndpoint) []string {
	return []string{
		endpoint.Region,
		string(endpoint.EndpointType),
		*endpoint.S3Endpoint,
	}
}
