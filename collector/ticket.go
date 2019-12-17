package collector

import (
	"context"
	"fmt"
	"log"

	"github.com/linode/linodego"
	"github.com/prometheus/client_golang/prometheus"
)

// TicketCollector represents a Linode Support Ticket
type TicketCollector struct {
	client linodego.Client

	Count *prometheus.Desc
}

// NewTicketCollector creates a TicketCollector
func NewTicketCollector(client linodego.Client) *TicketCollector {
	log.Println("[NewTicketCollector] Entered")
	fqName := name("tickets")
	labelKeys := []string{"status"}
	return &TicketCollector{
		client: client,

		Count: prometheus.NewDesc(
			fqName("count"),
			"Number of support tickets",
			labelKeys,
			nil,
		),
	}
}

// Collect implements Collector interface and is called by Prometheus to collect metrics
func (c *TicketCollector) Collect(ch chan<- prometheus.Metric) {
	log.Println("[TicketCollector:Collect] Entered")
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	tickets, err := c.client.ListTickets(ctx, nil)
	if err != nil {
		log.Println(err)
	}

	total := make(map[linodego.TicketStatus]float64)
	for _, t := range tickets {
		total[t.Status]++
	}
	for status, count := range total {
		ch <- prometheus.MustNewConstMetric(
			c.Count,
			prometheus.GaugeValue,
			count,
			// linodego.TicketStatus needs a String() method ;-)
			[]string{fmt.Sprintf("%s", status)}...,
		//[]string{status.String()}...,
		)
	}
	log.Println("[TicketCollector:Collect] Completes")
}

// Describe implements Collector interface and is called by Prometheus to describe metrics
func (c *TicketCollector) Describe(ch chan<- *prometheus.Desc) {
	log.Println("[TicketCollector:Describe] Entered")
	ch <- c.Count
	log.Println("[TicketCollector:Describe] Completes")
}
