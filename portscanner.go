package main

import (
	"fmt"
	"net"
	"strings"
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
	err    error
}

func (p *port) address() string {
	return fmt.Sprintf("%s:%d", p.host, p.number)
}

func (p *port) state() string {
	if p.err != nil {
		return "err"
	}

	if p.open {
		return "open"
	}

	return "closed"
}

func (ps *portscanner) test(p *port) {
	conn, err := net.DialTimeout("tcp", p.address(), ps.timeout)
	if conn != nil {
		defer conn.Close()
	}

	if err != nil {
		p.open = false
		if strings.Contains(err.Error(), "i/o timeout") != true && strings.Contains(err.Error(), "connection refused") != true {
			p.err = fmt.Errorf("error resolving tcp address for %s: %v\n", p.address(), err)
		}
	} else {
		p.open = true
	}
}

func newPortScanner(host string, timeout time.Duration) *portscanner {
	return &portscanner{host, timeout}
}

func (ps *portscanner) scan(concurrency, start, end int) (ports []*port) {
	fmt.Printf("scanning ports %d-%d on %s...\n", start, end, ps.host)

	var wg sync.WaitGroup
	sem := make(chan int, concurrency)
	openPortsChan := make(chan *port)

	for portNum := start; portNum <= end; portNum++ {
		p := &port{
			host:   ps.host,
			number: portNum,
			open:   false,
		}

		//ports = append(ports, p)
		wg.Add(1)
		go func(p *port) {
			sem <- 1

			go func(p *port) {
				defer wg.Done()

				ps.test(p)
				if p.open {
					openPortsChan <- p
				}
				<-sem
			}(p)
		}(p)
	}

	//closes channel
	go func() {
		wg.Wait()
		close(openPortsChan)
	}()

	for p := range openPortsChan {
		ports = append(ports, p)
	}

	return ports
}
