package collector

import (
	"context"
	"log"

	"github.com/linode/linodego"
	"github.com/prometheus/client_golang/prometheus"
)

// AccountCollector represents a Linode Account
type AccountCollector struct {
	client linodego.Client

	Balance *prometheus.Desc
	// Uninvoiced *prometheus.Desc
}

// NewAccountCollector creates an AccountCollector
func NewAccountCollector(client linodego.Client) *AccountCollector {
	log.Println("[NewAccountCollector] Entered")
	subsystem := "account"
	labelKeys := []string{"company", "email"}
	return &AccountCollector{
		client: client,

		Balance: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "balance"),
			"Balance of account",
			labelKeys,
			nil,
		),
		// Uninvoiced: prometheus.NewDesc(
		// 	prometheus.BuildFQName(namespace, subsystem, "uninvoiced"),
		// 	"Uninvoiced balance of account",
		// 	labelKeys,
		// 	nil,
		// ),
	}
}

// Collect implements Collector interface and is called by Prometheus to collect metrics
func (c *AccountCollector) Collect(ch chan<- prometheus.Metric) {
	log.Println("[AccountCollector:Collect] Entered")
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	account, err := c.client.GetAccount(ctx)
	if err != nil {
		//TODO(dazwilkin) capture logs
		log.Println(err)
	}

	ch <- prometheus.MustNewConstMetric(
		c.Balance,
		prometheus.GaugeValue,
		float64(account.Balance),
		[]string{account.Company, account.Email}...,
	)
	//TODO(dazwilkin) UnvoicedBalance is not yet implemented by the SDK
	// https://github.com/linode/linodego/issues/108
	// ch <- prometheus.MustNewConstMetric(
	// 	c.Uninvoiced,
	// 	prometheus.GaugeValue,
	// 	float64(account.BalanceUninvoiced),
	// 	[]string{account.Company, account.Email}...,
	// )
	log.Println("[AccountCollector:Collect] Completes")
}

// Describe implements Collector interface and is called by Prometheus to describe metrics
func (c *AccountCollector) Describe(ch chan<- *prometheus.Desc) {
	log.Println("[AccountCollector:Describe] Entered")
	ch <- c.Balance
	// ch <- c.Uninvoiced
	log.Println("[AccountCollector:Describe] Completes")
}
