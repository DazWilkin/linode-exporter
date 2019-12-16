package collector

import (
	"log"

	"github.com/linode/linodego"
	"github.com/prometheus/client_golang/prometheus"
)

type ManagedCollector struct {
	client linodego.Client

	CPUAvg    *prometheus.Desc
	CPUMax    *prometheus.Desc
	DiskAvg   *prometheus.Desc
	DiskMax   *prometheus.Desc
	SwapAvg   *prometheus.Desc
	SwapMax   *prometheus.Desc
	NetInAvg  *prometheus.Desc
	NetInMax  *prometheus.Desc
	NetOutAvg *prometheus.Desc
	NetOutMax *prometheus.Desc
}

func NewManagedCollector(client linodego.Client) *ManagedCollector {
	log.Println("[NewManagedCollector] Entered")
	fqName := name("managed")
	labelKeys := []string{"account"}
	return &ManagedCollector{
		client: client,
		CPUAvg: prometheus.NewDesc(
			fqName("cpu_average_utiliziation"),
			"CPU usage stats for the past 24 hours",
			labelKeys,
			nil,
		),
	}
}
func (c *ManagedCollector) Collect(ch chan<- prometheus.Metric) {
	log.Println("[ManagedCollector:Collect] Not yet implemented!")
	// ctx, cancel := context.WithTimeout(context.Background(), timeout)
	// defer cancel()

}
func (c *ManagedCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.CPUAvg
}
