package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"runtime"
	"time"

	"github.com/blockassets/bwpool_exporter/bwpool"
	"github.com/jpillora/overseer"
	"github.com/jpillora/overseer/fetcher"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	// Makefile build
	version = ""

	config  *string
	timeout *time.Duration
)

const (
	ghUser = "blockassets"
	ghRepo = "bwpool_exporter"
)

func main() {
	config = flag.String("config", "./bwpool.json", "Path to a file that has the API keys in it.")
	timeout = flag.Duration("timeout", 10*time.Second, "The amount of time to wait for bwpool to return.")
	port := flag.String("port", "5551", "The address to listen on for /metrics HTTP requests.")
	noUpdate := flag.Bool("no-update", false, "Never do any updates. Example: -no-update=true")

	flag.Parse()

	portStr := fmt.Sprintf(":%s", *port)

	if *noUpdate {
		prog(overseer.State{Address: portStr})
	} else {
		overseerRun(portStr, 1*time.Minute)
	}
}

func overseerRun(port string, interval time.Duration) {
	overseer.Run(overseer.Config{
		Program: prog,
		Address: port,
		Fetcher: &fetcher.Github{
			User:     ghUser,
			Repo:     ghRepo,
			Interval: interval,
		},
	})
}

func prog(state overseer.State) {
	log.Printf("%s %s %s %s on port %s\n", os.Args[0], version, runtime.GOOS, runtime.GOARCH, state.Address)

	poolConfig, err := bwpool.ReadConfig(*config)
	if err != nil {
		log.Fatal(err)
	}

	prometheus.MustRegister(NewExporter(poolConfig, *timeout))

	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.Serve(state.Listener, nil))
}
