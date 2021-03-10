package main

import (
	"GoPortScanner/src"
)

func main() {
	var (
		opts *src.ConfigOptions
		t    src.Target
		err  error
		c    chan bool
	)

	opts = new(src.ConfigOptions)

	opts, err = src.ParseFlags(opts)
	if err != nil {
		panic(err)
	}

	t.Ip = t.ParseIP(opts.Ip)
	t.Port = t.ParsePort(opts.Port)

	for _, i := range t.Ip {
		for _, p := range t.Port {
			t.Run(i, p, c)
		}
	}
}
