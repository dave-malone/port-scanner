package main

import (
	"flag"
	"fmt"
	"time"
)

var (
	timeout   = flag.Int("timeout", 1000, "port test timeout in milliseconds")
	hostname  = flag.String("hostname", "localhost", "Hostname or IP address of host to scan for open ports")
	startPort = flag.Int("start", 1, "Start port")
	endPort   = flag.Int("end", 49151, "End port")
	async     = flag.Bool("async", false, "Optionally run this process asynchronously. May produce unreliable results")
	maxConns  = flag.Int("max-conns", 250, "Maximum Concurrent Connections. Used only in async mode")
)

func main() {
	flag.Parse()

	t := time.Duration(*timeout) * time.Millisecond
	ps := newPortScanner(*hostname, t)

	start := time.Now()

	if *async {
		ps.ScanAsync(*maxConns, *startPort, *endPort)
	} else {
		ps.Scan(*startPort, *endPort)
	}

	elapsed := time.Since(start)
	fmt.Printf("Scan execution time: %s\n", elapsed)
}
