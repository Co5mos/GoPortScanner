package main

import (
	"log"
	"net"
	"os"
	"regexp"
	"strconv"
	"strings"
)

const (
	IPRegex = `\b(?:\d{1,3}\.){3}\d{1,3}\b$`
)

// 判断是否为IP地址
func IsIPAddr(ip string) bool {
	match, _ := regexp.MatchString(IPRegex, ip)
	return match
}

// C段转ip
func CidrToIPs(cidr string) []string {
	var ips []string

	// if a IP address, display the IP address and return
	if IsIPAddr(cidr) {
		log.Print("go run main.go 192.168.30.2 21 8080")
		os.Exit(1)
	}

	ipAddr, ipNet, err := net.ParseCIDR(cidr)
	if err != nil {
		log.Print("go run main.go 192.168.30.0/24 21 8080")
		os.Exit(1)
	}

	for ip := ipAddr.Mask(ipNet.Mask); ipNet.Contains(ip); increment(ip) {
		ips = append(ips, ip.String())
	}

	// CIDR too small eg. /31
	if len(ips) <= 2 {
		log.Print("go run main.go 192.168.30.0/24 21 8080")
		os.Exit(1)
	}

	return ips
}

// IP 段转 IP 列表
func Segment2IPs(segment string) []string {
	var ips []string

	// if a IP address, display the IP address and return
	if IsIPAddr(segment) {
		log.Print("go run main.go 192.168.30.2 21 8080")
		os.Exit(1)
	}

	ipRange := strings.Split(segment, "-")
	ipSegment := strings.Split(ipRange[0], ".")
	startNum, _ := strconv.Atoi(ipSegment[3])
	endNum, _ := strconv.Atoi(ipRange[1])

	for n := startNum; n <= endNum; n++ {
		ipTemp := append(ipSegment, strconv.Itoa(n))
		ip := strings.Join(ipTemp, ".")
		ips = append(ips, ip)
	}

	return ips
}
