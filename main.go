package main

type Target struct {
	ip   []string
	port []int
}

func main() {
	var opts *ConfigOptions
	var t Target
	var err error
	var c chan bool

	opts, err = ParseFlags(opts)
	if err != nil {
		panic(err)
	}

	t.ip = t.ParseIP(opts.ip)
	t.port = t.ParsePort(opts.port)

	for _, i := range t.ip {
		for _, p := range t.port {
			t.Run(i, p, c)
		}
	}
}
