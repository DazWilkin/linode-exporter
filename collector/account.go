package collector

import (
	"context"
	"log"

	"github.com/linode/linodego"
	"github.com/prometheus/client_golang/prometheus"
)

type AccountCollector struct {
	client linodego.Client

	Balance    *prometheus.Desc
	Uninvoiced *prometheus.Desc
}

func NewAccountCollector(client linodego.Client) *AccountCollector {
	labelKeys := []string{"company", "email"}
	return &AccountCollector{
		client: client,

		Balance: prometheus.NewDesc(
			prefix+"_account_balance",
			"Balance of account",
			labelKeys,
			nil,
		),
		Uninvoiced: prometheus.NewDesc(
			prefix+"_account_uninvoiced",
			"Uninvoiced balance of account",
			labelKeys,
			nil,
		),
	}
}
func (c *AccountCollector) Collect(ch chan<- prometheus.Metric) {
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
	ch <- prometheus.MustNewConstMetric(
		c.Uninvoiced,
		prometheus.GaugeValue,
		0.0, //float64(account.BalanceUninvoiced),
		[]string{account.Company, account.Email}...,
	)

}
func (c *AccountCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.Balance
	ch <- c.Uninvoiced
}
