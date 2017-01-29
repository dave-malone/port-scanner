package main

import (
	"fmt"
	"net"
	"sync"
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
	//fmt.Printf("testing %s w/ timeout %v\n", p.address(), timeout)

	tcpAddr, err := net.ResolveTCPAddr("tcp4", p.address())
	if err != nil {
		p.open = false
		return fmt.Errorf("error resolving tcp4 address for %s: %v", p.address(), err)
	}

	conn, err := net.DialTimeout("tcp", tcpAddr.String(), timeout)
	if conn != nil {
		defer conn.Close()
	}

	if err != nil {
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

	var wg sync.WaitGroup

	for portNum := start; portNum <= end; portNum++ {
		p := port{
			host:   ps.host,
			number: portNum,
			open:   false,
		}

		wg.Add(1)

		go func(p port) {
			// Decrement the counter when the goroutine completes.
			defer wg.Done()

			if err := p.test(ps.timeout); err != nil {
				fmt.Printf("%s [Failed to test port: %v]\n", p.address(), err)
				return
			}

			if p.open {
				fmt.Printf("%s [%s]\n", p.address(), p.portState())
			}
		}(p)
	}

	wg.Wait()
}
