package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"text/template"
	"time"

	"github.com/DazWilkin/linode-exporter/collector"

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
	token       = flag.String("linode_token", os.Getenv("LINODE_TOKEN"), "Linode API Token")
	debug       = flag.Bool("debug", false, "Enable Linode REST API debugging")
	endpoint    = flag.String("endpoint", ":9388", "The endpoint of the HTTP server")
	metricsPath = flag.String("path", "/metrics", "The path on which Prometheus metrics will be served")
)

const (
	rootContent = `<!DOCTYPE html>
<html>
	<head>
		<title>Linode Exporter</title>
	</head>
	<body>
		<h1>Linode Exporter</h1>
		<p><a href="{{.MetricsPath}}">Metrics</a></p>
	</body>
</html>`
)

var (
	rootTemplate = template.Must(template.New("root").Parse(rootContent))
)

func rootHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=UTF-8")
	if err := rootTemplate.Execute(w, struct {
		MetricsPath string
	}{
		MetricsPath: *metricsPath,
	}); err != nil {
		log.Printf("[rootHandler] error executing template: %v", err)
	}
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
	registry.MustRegister(collector.NewAccountCollector(client))
	registry.MustRegister(collector.NewExporterCollector(client, OSVersion, GitCommit))
	registry.MustRegister(collector.NewInstanceCollector(client))
	registry.MustRegister(collector.NewInstanceStatsCollector(client))
	registry.MustRegister(collector.NewKubernetesCollector(client))
	registry.MustRegister(collector.NewNodeBalancerCollector(client))
	registry.MustRegister(collector.NewTicketCollector(client))
	registry.MustRegister(collector.NewVolumeCollector(client))
	registry.MustRegister(collector.NewObjectStorageCollector(client))

	mux := http.NewServeMux()
	mux.Handle("/", http.HandlerFunc(rootHandler))
	mux.Handle(*metricsPath, promhttp.HandlerFor(registry, promhttp.HandlerOpts{
		Timeout: timeout,
	}))

	log.Printf("[main] Server starting (%s)", *endpoint)
	log.Printf("[main] metrics served on: %s", *metricsPath)
	log.Fatal(http.ListenAndServe(*endpoint, mux))
}
