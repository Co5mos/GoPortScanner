package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	browser "github.com/EDDYCJY/fake-useragent"
	"github.com/antchfx/htmlquery"
)

const (
	IPRegex = `\b(?:\d{1,3}\.){3}\d{1,3}\b$`
)

var (
	h         bool   // Show help
	ipAddr    string // IP地址
	cidrAddr  string // C段
	ipSeg     string // IP段
	startPort int    // 开始端口
	endPort   int    // 结束端口
)

func init() {
	flag.BoolVar(&h, "h", false, "Show help")
	flag.StringVar(&ipAddr, "i", "127.0.0.1", "IP地址")
	flag.StringVar(&cidrAddr, "c", "192.168.30.0/24", "C段地址")
	flag.StringVar(&ipSeg, "ips", "192.168.30.1-20", "IP段")
	flag.IntVar(&startPort, "s", 21, "开始端口")
	flag.IntVar(&endPort, "e", 9000, "结束端口")
}

func getPorts(startPort int, endPort int) []int {
	// 获取端口列表
	var ports []int

	for i := startPort; i <= endPort; i++ {
		ports = append(ports, i)
	}

	return ports
}

func cidr2IPs(cidr string) []string {
	// C段转ip
	var ips []string

	// if a IP address, display the IP address and return
	if isIPAddr(cidr) {
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

func segment2IPs(segment string) []string {
	// ip段转ip
	var ips []string

	// if a IP address, display the IP address and return
	if isIPAddr(segment) {
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

func scanPorst(ip string, port int) bool {
	// 扫描端口
	ipPort := fmt.Sprintf("%s:%d", ip, port)

	conn, err := net.DialTimeout("tcp", ipPort, 200*time.Millisecond)
	if err != nil {
		log.Print(err)
		return false
	}

	conn.Close()

	return true
}

func getHTTPBanner(url string) {
	// 获取网站banner
	randomUA := browser.Random()

	client := &http.Client{}
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("User-Agent", randomUA)

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("Http get err:", err)
	}
	defer resp.Body.Close()

	doc, err := htmlquery.Parse(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	titles, err := htmlquery.QueryAll(doc, `/html/head/title/text()`)
	if err != nil {
		panic(`not a valid XPath expression.`)
	}

	for _, title := range titles {
		fmt.Printf("[+] %s: %s\n", url, htmlquery.InnerText(title))
	}

	// fmt.Printf("[+] %s: %s\n", url, htmlquery.InnerText(title))
	// return htmlquery.InnerText(title)
}

func getScoketInfo(ip string, port int) {
	// 获取socket信息
	ipPort := fmt.Sprintf(ip, ":", port)

	_, err := net.DialTimeout("tcp", ipPort, 200*time.Millisecond)
	if err != nil {
		log.Print(err)
		return
	}

}

func isIPAddr(cidr string) bool {
	match, _ := regexp.MatchString(IPRegex, cidr)
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

func run(ip string, port int, ch chan bool) {
	if scanPorst(ip, port) {
		url := fmt.Sprintf("http://%s:%d", ip, port)
		getHTTPBanner(url)
	}
	ch <- true
}

func checkTarget(target string, startPort int, endPort int) {
	// 检测目标
	// 生成port切片
	ch := make(chan bool)
	ports := getPorts(startPort, endPort)
	ipSize := 1

	if strings.Contains(target, "/") {
		// C段查询
		ipSlice := cidr2IPs(target)
		ipSize = len(ipSlice) - 2
		for _, ip := range ipSlice[1 : len(ipSlice)-1] {
			for _, port := range ports {
				go run(ip, port, ch)
			}
		}

	} else if strings.Contains(target, "-") {
		// IP段查询
		ipSlice := segment2IPs(target)
		for _, ip := range ipSlice[1:len(ipSlice)] {
			fmt.Println(ip)
		}
	} else {
		// 单个IP
		for _, port := range ports[0:len(ports)] {
			go run(target, port, ch)
		}
	}

	scanSize := ipSize * len(ports) // 扫描的所有IP*Port

	for i := 0; i < scanSize; i++ {
		<-ch
	}
}

func main() {
	startPort, _ := strconv.Atoi(os.Args[2])
	endPort, _ := strconv.Atoi(os.Args[3])
	checkTarget(os.Args[1], startPort, endPort)
}
