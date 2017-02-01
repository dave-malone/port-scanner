package main

import (
	"fmt"
	"net"
	"os/exec"
)

type host struct {
	ip        net.IP
	alive     bool
	names     []string
	openPorts []port
}

func (h *host) lookupNames() error {
	names, err := net.LookupAddr(h.ip.String())
	h.names = names
	return err
}

func newHost(hostname string) *host {
	ip := net.ParseIP(hostname)
	if ip == nil {
		tcpAddr, err := net.ResolveTCPAddr("tcp4", fmt.Sprintf("%s:%d", hostname, 80))
		if err == nil && tcpAddr != nil {
			ip = tcpAddr.IP
		}
	}

	if ip == nil {
		panic("Failed to initialize new host by hostname " + hostname)
	}

	host := &host{
		ip: ip,
	}
	host.lookupNames()

	return host
}

func mapNetwork(cidr string) ([]host, error) {
	ipnet, ips, err := listIPsInCIDR(cidr)
	if err != nil {
		return nil, err
	}

	fmt.Printf("searching for active hosts on %s\n", ipnet)

	concurrentMax := 100
	pingChan := make(chan net.IP, concurrentMax)
	hostChan := make(chan host, len(ips))
	doneChan := make(chan []host)

	for i := 0; i < concurrentMax; i++ {
		go ping(pingChan, hostChan)
	}

	go receivePong(len(ips), hostChan, doneChan)

	for _, ip := range ips {
		pingChan <- ip
	}

	alives := <-doneChan

	return alives, nil
}

func ping(pingChan <-chan net.IP, hostChan chan<- host) {
	for ip := range pingChan {
		_, err := exec.Command("ping", "-c1", "-t1", ip.String()).Output()
		var alive bool
		if err != nil {
			alive = false
		} else {
			alive = true
		}
		hostChan <- host{ip: ip, alive: alive}
	}
}

func receivePong(hostNum int, hostChan <-chan host, doneChan chan<- []host) {
	var alives []host
	for i := 0; i < hostNum; i++ {
		host := <-hostChan
		if host.alive {
			host.lookupNames()
			alives = append(alives, host)
		}
	}
	doneChan <- alives
}

func listIPsInCIDR(cidr string) (*net.IPNet, []net.IP, error) {
	baseip, ipnet, err := net.ParseCIDR(cidr)
	if err != nil {
		return nil, nil, fmt.Errorf("Failed to parse cidr %s; %v", cidr, err)
	}

	ips := make([]net.IP, 0)
	for ip := baseip.Mask(ipnet.Mask); ipnet.Contains(ip); ip = inc(ip) {
		ips = append(ips, ip)
	}

	// remove network address and broadcast address
	return ipnet, ips[1 : len(ips)-1], nil
}

func inc(ip net.IP) net.IP {
	nextip := make(net.IP, len(ip))
	copy(nextip, ip)
	for j := len(nextip) - 1; j >= 0; j-- {
		nextip[j]++
		if nextip[j] > 0 {
			break
		}
	}

	return nextip
}
