package collector

import (
	"log"
	"runtime"

	"github.com/linode/linodego"
	"github.com/prometheus/client_golang/prometheus"
)

// ExporterCollector represents the Linode Exporter (itself)
type ExporterCollector struct {
	client linodego.Client

	Up *prometheus.Desc

	osVersion string
	gitCommit string
}

// NewExporterCollector creates an ExporterCollector
func NewExporterCollector(client linodego.Client, osVersion, gitCommit string) *ExporterCollector {
	log.Println("[NewExporterCollector] Entered")
	subsystem := "exporter"
	labelKeys := []string{"goVersion", "osVersion", "exporterCommit"}
	return &ExporterCollector{
		client: client,
		Up: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "info"),
			"A metric with a constant value of '1' labeled with go, OS and the exporter versions",
			labelKeys,
			nil,
		),
		osVersion: func() string {
			log.Printf("[NewExporterCollector] OS=%s", osVersion)
			return osVersion
		}(),
		gitCommit: func() string {
			log.Printf("[NewExporterCollector] Commit=%s", gitCommit)
			return gitCommit
		}(),
	}
}

// Collect implements Collector interface and is called by Prometheus to collect metrics
func (c *ExporterCollector) Collect(ch chan<- prometheus.Metric) {
	log.Println("[ExporterCollector:Collect] Entered")
	ch <- prometheus.MustNewConstMetric(
		c.Up,
		prometheus.CounterValue,
		1,
		[]string{
			runtime.Version(),
			c.osVersion,
			c.gitCommit,
		}...)
	log.Println("[ExporterCollector:Collect] Completes")
}

// Describe implements Collector interface and is called by Prometheus to describe metrics
func (c *ExporterCollector) Describe(ch chan<- *prometheus.Desc) {
	log.Println("[ExporterCollector:Describe] Entered")
	ch <- c.Up
	log.Println("[ExporterCollector:Describe] Completes")
}
