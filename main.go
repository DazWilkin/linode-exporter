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
	apiToken    = flag.String("token", "", "Linode API Token")
	endpoint    = flag.String("endpoint", ":2112", "The endpoint of the HTTP server")
	metricsPath = flag.String("path", "/metrics", "The path on which Prometheus metrics will be served")
)

func main() {
	flag.Parse()
	if *apiToken == "" {
		log.Fatal("Provide Linode API Token")
	}
	source := oauth2.StaticTokenSource(&oauth2.Token{
		AccessToken: *apiToken,
	})
	oauth2Client := &http.Client{
		Transport: &oauth2.Transport{
			Source: source,
		},
	}
	linodeClient := linodego.NewClient(oauth2Client)
	linodeClient.SetDebug(true)

	registry := prometheus.NewRegistry()
	registry.MustRegister(collector.NewInstanceCollector(linodeClient))
	registry.MustRegister(collector.NewNodeBalancerCollector(linodeClient))

	http.Handle(*metricsPath, promhttp.HandlerFor(registry, promhttp.HandlerOpts{}))
	log.Fatal(http.ListenAndServe(*endpoint, nil))
}
