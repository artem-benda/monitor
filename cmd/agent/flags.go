package main

import "flag"

var (
	serverEndpoint string
	reportInterval int
	pollInterval   int
)

func parseFlags() {
	flag.StringVar(&serverEndpoint, "a", "localhost:8080", "address and port of metrics server")
	flag.IntVar(&reportInterval, "r", 10, "send metrics delay in seconds")
	flag.IntVar(&pollInterval, "p", 2, "poll runtime metrics delay in seconds")
	flag.Parse()
}
