package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/DazWilkin/linode-exporter/collector"

	"github.com/linode/linodego"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"golang.org/x/oauth2"
)

var (
	token       = flag.String("linode_token", "", "Linode API Token")
	debug       = flag.Bool("debug", false, "Enable Linode REST API debugging")
	endpoint    = flag.String("endpoint", ":9388", "The endpoint of the HTTP server")
	metricsPath = flag.String("path", "/metrics", "The path on which Prometheus metrics will be served")
)

func main() {
	flag.Parse()
	if *token == "" {
		log.Fatal("Provide Linode API Token")
	}
	source := oauth2.StaticTokenSource(&oauth2.Token{
		AccessToken: *token,
	})
	oauth2Client := &http.Client{
		Transport: &oauth2.Transport{
			Source: source,
		},
	}
	client := linodego.NewClient(oauth2Client)
	client.SetDebug(*debug)

	registry := prometheus.NewRegistry()
	registry.MustRegister(collector.NewAccountCollector(client))
	registry.MustRegister(collector.NewInstanceCollector(client))
	registry.MustRegister(collector.NewNodeBalancerCollector(client))
	registry.MustRegister(collector.NewTicketCollector(client))

	http.Handle(*metricsPath, promhttp.HandlerFor(registry, promhttp.HandlerOpts{}))
	log.Fatal(http.ListenAndServe(*endpoint, nil))
}
