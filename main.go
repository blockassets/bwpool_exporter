package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/blockassets/bwpool_exporter/bwpool"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	// Makefile build
	version = ""
)

func main() {
	config := flag.String("config", "./bwpool.json", "Path to a file that has the API keys in it.")
	port := flag.String("port", "5551", "The address to listen on for /metrics HTTP requests.")
	timeout := flag.Duration("timeout", 10*time.Second, "The amount of time to wait for bwpool to return.")
	flag.Parse()

	poolConfig, err := bwpool.ReadConfig(*config)
	if err != nil {
		log.Fatal(err)
	}

	prometheus.MustRegister(NewExporter(poolConfig, *timeout))

	http.Handle("/metrics", promhttp.Handler())
	log.Printf("%s %s", os.Args[0], version)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", *port), nil))
}
