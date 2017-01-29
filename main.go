package main

import (
	"flag"
	"time"
)

var (
	timeout   = flag.Int("timeout", 1000, "port test timeout in milliseconds")
	hostname  = flag.String("hostname", "localhost", "Hostname or IP address of host to scan for open ports")
	startPort = flag.Int("start", 1, "Start port")
	endPort   = flag.Int("end", 100, "End port")
)

func main() {
	flag.Parse()

	t := time.Duration(*timeout) * time.Millisecond
	ps := newPortScanner(*hostname, t)

	ps.Scan(*startPort, *endPort)
}
