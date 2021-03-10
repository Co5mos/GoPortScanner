package main

import (
	"strconv"
	"strings"
)

// 解析 Port 参数
func (t *Target) ParsePort(port string) []int {
	portSegment := strings.Split(port, ",")

	// 解析端口范围 80-90,999
	for _, ps := range portSegment {
		p, err := strconv.Atoi(ps)
		if err == nil {
			// 单个 port 直接加入切片
			t.port = append(t.port, p)

		} else if strings.Contains(ps, "-") {
			// 80-90
			pp := strings.Split(ps, "-")

			// 解析错误，返回空切片
			if len(pp) < 2 {
				return t.port
			}

			pstart, err := strconv.Atoi(pp[0])
			if err != nil {
				panic(err)
			}

			pend, err := strconv.Atoi(pp[1])
			if err != nil {
				panic(err)
			}

			for i := pstart; i <= pend; i++ {
				t.port = append(t.port, i)
			}
		}
	}

	return t.port
}
