package main

import (
	"errors"
	"flag"
	"fmt"
)

// 解析校验参数

type ConfigOptions struct {
	ip   string
	port string
}

func Usage() {

}

// 解析命令行参数
func ParseFlags(opts *ConfigOptions) (*ConfigOptions, error) {
	flag.StringVar(&opts.ip, "i", "127.0.0.1", "IP地址")
	//flag.StringVar(&opts.IP.cidr, "c", "192.168.30.0/24", "C段地址")
	//flag.StringVar(&opts.IP.seg, "ips", "192.168.30.1-20", "IP段")

	flag.StringVar(&opts.port, "p", "80", "端口号")
	//flag.IntVar(&opts.Port.startPort, "s", 21, "开始端口")
	//flag.IntVar(&opts.Port.endPort, "e", 9000, "结束端口")

	flag.Usage = Usage
	flag.Parse()

	if flag.NArg() != 0 {
		flag.Usage()
		return nil, errors.New(fmt.Sprintf("Parse options error, please see usage for more information."))
	}

	return opts, nil
}
