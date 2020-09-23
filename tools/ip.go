package tools

import (
	"net"
)

func ExternalIP() (ips []*net.IPNet, err error) {
	interfaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}
	for _, v := range interfaces {
		if v.Flags&net.FlagUp != net.FlagUp {
			continue
		}
		if v.Flags&net.FlagLoopback == net.FlagLoopback {
			continue
		}
		addrs, err := v.Addrs()
		if err != nil {
			continue
		}
		for _, addr := range addrs {
			a, ok := addr.(*net.IPNet)
			if ok {
				ips = append(ips, a)
			}
		}
	}
	return
}
