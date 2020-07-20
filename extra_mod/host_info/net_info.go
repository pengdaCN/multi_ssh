package host_info

import (
	"github.com/tidwall/gjson"
	"net"
	"strings"
)

type NetInterInfo struct {
	net.Interface
	Ip []net.IPNet
}

func ParseNetInterInfoWithJsonStr(jstr string) (*NetInterInfo, error) {
	j := gjson.Parse(jstr)
	return ParseNetInterInfoWithGjson(j)
}

func ParseNetInterInfoWithGjson(inter gjson.Result) (*NetInterInfo, error) {
	var hw net.Interface
	hw.Name = inter.Get("ifname").String()
	hw.Index = int(inter.Get("ifindex").Int())
	hw.MTU = int(inter.Get("mtu").Int())
	m, err := net.ParseMAC(inter.Get("address").String())
	if err != nil {
		return nil, err
	}
	hw.HardwareAddr = m
	flags := inter.Get("flags").Array()
	for i := 0; i < len(flags); i++ {
		if flags[i].String() == "LOWER_UP" {
			hw.Flags = net.FlagUp
			break
		}
	}
	ips := inter.Get("addr_info").Array()
	ipArr := make([]net.IPNet, 0)
	for i := 0; i < len(ips); i++ {
		ip := net.ParseIP(ips[i].Get("local").String())
		var mask net.IPMask
		if strings.Contains(ip.String(), ".") {
			mask = net.CIDRMask(int(ips[i].Get("prefixlen").Int()), 32)
		} else {
			mask = net.CIDRMask(int(ips[i].Get("prefixlen").Int()), 128)
		}
		ipnet := net.IPNet{
			IP:   ip,
			Mask: mask,
		}
		ipArr = append(ipArr, ipnet)
	}
	return &NetInterInfo{
		Interface: hw,
		Ip:        ipArr,
	}, nil
}

func (n *NetInterInfo) Active() bool {
	if n.Flags != net.FlagUp {
		return false
	}
	return true
}

func (n *NetInterInfo) Ipv4() *net.IPNet {
	for _, i := range n.Ip {
		if ip := i.IP.To4(); ip != nil {
			return &i
		}
	}
	return nil
}

func (n *NetInterInfo) Ipv4Str() string {
	var ipv4Str string
	if i := n.Ipv4(); i != nil {
		ipv4Str = i.String()
	}
	return ipv4Str
}

func (n *NetInterInfo) Ipv6() *net.IPNet {
	for _, i := range n.Ip {
		if ip := i.IP.To16(); ip != nil {
			return &i
		}
	}
	return nil
}

func (n *NetInterInfo) Ipv6Str() string {
	var ipv6Str string
	if i := n.Ipv6(); i != nil {
		ipv6Str = i.String()
	}
	return ipv6Str
}
