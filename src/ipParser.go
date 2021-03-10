package src

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
func (t *Target) IsIPAddr(ip string) bool {
	match, _ := regexp.MatchString(IPRegex, ip)
	return match
}

func increment(ip net.IP) {
	for i := len(ip) - 1; i >= 0; i-- {
		ip[i]++
		if ip[i] != 0 {
			break
		}
	}
}

// C段转ip
func (t *Target) CidrToIPs(cidr string) []string {
	var ips []string

	// if a IP address, display the IP address and return
	if t.IsIPAddr(cidr) {
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
func (t *Target) SegmentToIPs(segment string) []string {
	var ips []string

	// if a IP address, display the IP address and return
	if t.IsIPAddr(segment) {
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

// 解析 IP 参数
func (t *Target) ParseIP(ip string) []string {
	if strings.Contains(ip, "/") {
		// C段查询
		ipSlice := t.CidrToIPs(ip)
		for _, ip := range ipSlice[1 : len(ipSlice)-1] {
			t.Ip = append(t.Ip, ip)
		}

	} else if strings.Contains(ip, "-") {
		// IP段查询
		ipSlice := t.SegmentToIPs(ip)
		for _, ip := range ipSlice[1 : len(ipSlice)-1] {
			t.Ip = append(t.Ip, ip)
		}

	} else {
		// 单个IP
		t.Ip = append(t.Ip, ip)
	}

	return t.Ip
}
