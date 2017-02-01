package main

import (
	"fmt"
	"net"
	"os/exec"
)

type pong struct {
	ip    net.IP
	alive bool
}

func mapnetwork(cidr string) ([]net.IP, error) {
	ipnet, ips, err := getallips(cidr)
	if err != nil {
		return nil, err
	}

	fmt.Printf("searching for active hosts on %s\n", ipnet)

	concurrentMax := 100
	pingChan := make(chan net.IP, concurrentMax)
	pongChan := make(chan pong, len(ips))
	doneChan := make(chan []net.IP)

	for i := 0; i < concurrentMax; i++ {
		go ping(pingChan, pongChan)
	}

	go receivePong(len(ips), pongChan, doneChan)

	for _, ip := range ips {
		pingChan <- ip
		//  fmt.Println("sent: " + ip)
	}

	alives := <-doneChan

	return alives, nil
}

func ping(pingChan <-chan net.IP, pongChan chan<- pong) {
	for ip := range pingChan {
		_, err := exec.Command("ping", "-c1", "-t1", ip.String()).Output()
		var alive bool
		if err != nil {
			alive = false
		} else {
			alive = true
		}
		pongChan <- pong{ip: ip, alive: alive}
	}
}

func receivePong(pongNum int, pongChan <-chan pong, doneChan chan<- []net.IP) {
	var alives []net.IP
	for i := 0; i < pongNum; i++ {
		pong := <-pongChan
		//  fmt.Println("received:", pong)
		if pong.alive {
			alives = append(alives, pong.ip)
		}
	}
	doneChan <- alives
}

func getallips(cidr string) (*net.IPNet, []net.IP, error) {
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
