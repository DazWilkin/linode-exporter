package main

import (
	"flag"
	"fmt"
	"log"
	"maps"
	"net/http"
	"os"
	"slices"
	"strings"
	"time"

	"main/collector"

	"github.com/linode/linodego"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"golang.org/x/oauth2"
)

const (
	timeout = 60 * time.Second
)

var (
	// GitCommit expected to be set during build and contain the git commit value
	GitCommit string
	// OSVersion expected to be set during build and contain the OS version (uname --kernel-release)
	OSVersion string
)
var (
	token              = flag.String("linode-token", os.Getenv("LINODE_TOKEN"), "Linode API Token")
	debug              = flag.Bool("debug", false, "Enable Linode REST API debugging")
	endpoint           = flag.String("endpoint", ":9388", "The endpoint of the HTTP server")
	metricsPath        = flag.String("path", "/metrics", "The path on which Prometheus metrics will be served")
	enabledCollectors  = flag.String("collectors", os.Getenv("COLLECTORS"), "Comma-separated list of enabled collectors (default: all)")
	collectorFactories = map[string]func(linodego.Client) prometheus.Collector{
		"account": func(c linodego.Client) prometheus.Collector { return collector.NewAccountCollector(c) },
		"exporter": func(c linodego.Client) prometheus.Collector {
			return collector.NewExporterCollector(c, OSVersion, GitCommit)
		},
		"instance":       func(c linodego.Client) prometheus.Collector { return collector.NewInstanceCollector(c) },
		"instance_stats": func(c linodego.Client) prometheus.Collector { return collector.NewInstanceStatsCollector(c) },
		"kubernetes":     func(c linodego.Client) prometheus.Collector { return collector.NewKubernetesCollector(c) },
		"nodebalancer":   func(c linodego.Client) prometheus.Collector { return collector.NewNodeBalancerCollector(c) },
		"ticket":         func(c linodego.Client) prometheus.Collector { return collector.NewTicketCollector(c) },
		"volume":         func(c linodego.Client) prometheus.Collector { return collector.NewVolumeCollector(c) },
		"objectstorage":  func(c linodego.Client) prometheus.Collector { return collector.NewObjectStorageCollector(c) },
	}
)

func enableCollectors(client linodego.Client, registry *prometheus.Registry) {
	// If no environment variable is set, enable all collectors
	if *enabledCollectors == "" {
		*enabledCollectors = strings.Join(slices.Collect(maps.Keys(collectorFactories)), ",")
	}
	for _, name := range strings.Split(*enabledCollectors, ",") {
		name = strings.TrimSpace(name)
		log.Printf("[getEnabledCollectors] Creating collector %s", name)
		if factory, exists := collectorFactories[name]; exists {
			registry.MustRegister(factory(client))
		} else {
			log.Printf("[getEnabledCollectors] Collector %s not found", name)
		}
	}
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=UTF-8")
	fmt.Fprint(w, "<h2>Linode Exporter</h2>")
	fmt.Fprintf(w, "<a href=\"%s\">metrics</a>", *metricsPath)
}
func main() {
	flag.Parse()
	if *token == "" {
		log.Fatal("Provide Linode API Token")
	}

	if GitCommit == "" {
		log.Println("[main] GitCommit value unset (\"\"); expected to be set during build")
	}
	if OSVersion == "" {
		log.Println("[main] OSVersion value (\"\"); expected to be set during build")
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
	enableCollectors(client, registry)

	mux := http.NewServeMux()
	mux.Handle("/", http.HandlerFunc(rootHandler))
	mux.Handle(*metricsPath, promhttp.HandlerFor(registry, promhttp.HandlerOpts{
		Timeout: timeout,
	}))

	log.Printf("[main] Server starting (%s)", *endpoint)
	log.Printf("[main] metrics served on: %s", *metricsPath)
	log.Fatal(http.ListenAndServe(*endpoint, mux))
}
