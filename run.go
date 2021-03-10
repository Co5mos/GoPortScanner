package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"time"

	UA "github.com/EDDYCJY/fake-useragent"
	"github.com/antchfx/htmlquery"
)

func (t *Target) ScanPort(ip string, port int) bool {
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

func (t *Target) GetHTTPBanner(url string) {
	// 获取网站banner
	randomUA := UA.Random()

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

func (t *Target) GetScoketInfo(ip string, port int) {
	// 获取socket信息
	ipPort := fmt.Sprintf(ip, ":", port)

	_, err := net.DialTimeout("tcp", ipPort, 200*time.Millisecond)
	if err != nil {
		log.Print(err)
		return
	}

}

func (t *Target) Run(ip string, port int, ch chan bool) {
	if t.ScanPort(ip, port) {
		url := fmt.Sprintf("http://%s:%d", ip, port)
		t.GetHTTPBanner(url)
	}
	ch <- true
}
