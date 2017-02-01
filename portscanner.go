package main

import (
	"fmt"
	"net"
	"strings"
	"time"
)

type portscanner struct {
	host    string
	timeout time.Duration
}

type port struct {
	host   string
	number int
	open   bool
}

func (p *port) address() string {
	return fmt.Sprintf("%s:%d", p.host, p.number)
}

func (p *port) portState() string {
	if p.open {
		return "open"
	}

	return "closed"
}

func (p *port) test(timeout time.Duration) error {
	conn, err := net.DialTimeout("tcp", p.address(), timeout)
	if conn != nil {
		defer conn.Close()
	}

	if err != nil {
		errMsg := err.Error()
		if strings.Contains(errMsg, "i/o timeout") != true && strings.Contains(errMsg, "connection refused") != true {
			fmt.Printf("error resolving tcp4 address for %s: %v\n", p.address(), err)
		}

		p.open = false
		return nil
	}

	p.open = true
	return nil
}

func newPortScanner(host string, timeout time.Duration) *portscanner {
	return &portscanner{host, timeout}
}

func (ps *portscanner) Scan(start int, end int) {
	fmt.Printf("scanning ports %d-%d...\n", start, end)

	for portNum := start; portNum <= end; portNum++ {
		p := port{
			host:   ps.host,
			number: portNum,
			open:   false,
		}

		if err := p.test(ps.timeout); err != nil {
			fmt.Printf("%s [Failed to test port: %v]\n", p.address(), err)
			return
		}

		if p.open {
			fmt.Printf("%s [%s]\n", p.address(), p.portState())
		}
	}
}

//running the scan async does not seem to produce reliable results
func (ps *portscanner) ScanAsync(maxConns, start, end int) {
	fmt.Printf("[async] scanning ports %d-%d...\n", start, end)

	//var wg sync.WaitGroup

	concurrency := maxConns
	sem := make(chan bool, concurrency)

	for portNum := start; portNum <= end; portNum++ {
		p := port{
			host:   ps.host,
			number: portNum,
			open:   false,
		}

		sem <- true

		go func(p port) {
			defer func() { <-sem }()

			if err := p.test(ps.timeout); err != nil {
				fmt.Printf("%s [Failed to test port: %v]\n", p.address(), err)
				return
			}

			if p.open {
				fmt.Printf("%s [%s]\n", p.address(), p.portState())
			}
		}(p)
	}

	for i := 0; i < cap(sem); i++ {
		sem <- true
	}
}
