package main

import (
	"flag"
	"fmt"
	"net"
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
	// bytetest()
	nmap()
}

func bytetest() {
	var b byte

	fmt.Printf("%08b\n", b)
	for b++; b > 0; b++ {
		fmt.Printf("%08b\n", b)
	}
}

func nmap() {
	cidr := "192.168.1.1/24"
	alives, err := mapnetwork(cidr)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	for _, ip := range alives {
		fmt.Printf("Scanning %s:", ip)
		if names, err := net.LookupAddr(ip.String()); err == nil {
			fmt.Printf("%v", names)
		}
		fmt.Println(" for open ports:")
		portscan(ip.String())
	}
}

func portscan(host string) {
	t := time.Duration(*timeout) * time.Millisecond
	ps := newPortScanner(host, t)

	start := time.Now()

	if *async {
		ps.ScanAsync(*maxConns, *startPort, *endPort)
	} else {
		ps.Scan(*startPort, *endPort)
	}

	elapsed := time.Since(start)
	fmt.Printf("Scan execution time: %s\n", elapsed)
}
