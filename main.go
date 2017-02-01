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
	maxConns  = flag.Int("max-conns", 1000, "Maximum Concurrent Connections.")
)

func main() {
	flag.Parse()
	// bytetest()
	//scanAndReportOpenPorts(*newHost(*hostname))
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
	alives, err := mapNetwork(cidr)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	for _, host := range alives {
		scanAndReportOpenPorts(host)
	}
}

func scanAndReportOpenPorts(host host) {
	fmt.Printf("Scanning %s:%v for open ports:\n", host.ip, host.names)
	start := time.Now()
	ports := portscan(host.ip.String())
	elapsed := time.Since(start)
	for _, p := range ports {
		if p.open {
			fmt.Printf("%d [%s]\n", p.number, p.state())
		}
	}
	fmt.Printf("Scan execution time: %s\n", elapsed)
	fmt.Println()
}

func portscan(host string) []*port {
	t := time.Duration(*timeout) * time.Millisecond
	ps := newPortScanner(host, t)
	return ps.scan(*maxConns, *startPort, *endPort)
}
