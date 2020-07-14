package main

import (
	"fmt"
	"multi_ssh/m_terminal"
	"multi_ssh/model"
	"regexp"
	"strings"
	"testing"
)

var str = `[
  {
    "ifindex": 1,
    "ifname": "lo",
    "flags": [
      "LOOPBACK",
      "UP",
      "LOWER_UP"
    ],
    "mtu": 65536,
    "qdisc": "noqueue",
    "operstate": "UNKNOWN",
    "group": "default",
    "txqlen": 1000,
    "link_type": "loopback",
    "address": "00:00:00:00:00:00",
    "broadcast": "00:00:00:00:00:00",
    "addr_info": [
      {
        "family": "inet",
        "local": "127.0.0.1",
        "prefixlen": 8,
        "scope": "host",
        "label": "lo",
        "valid_life_time": 4294967295,
        "preferred_life_time": 4294967295
      },
      {
        "family": "inet6",
        "local": "::1",
        "prefixlen": 128,
        "scope": "host",
        "valid_life_time": 4294967295,
        "preferred_life_time": 4294967295
      }
    ]
  },
  {
    "ifindex": 2,
    "ifname": "enp3s0",
    "flags": [
      "NO-CARRIER",
      "BROADCAST",
      "MULTICAST",
      "UP"
    ],
    "mtu": 1500,
    "qdisc": "fq_codel",
    "operstate": "DOWN",
    "group": "default",
    "txqlen": 1000,
    "link_type": "ether",
    "address": "d8:c4:97:92:f8:b0",
    "broadcast": "ff:ff:ff:ff:ff:ff",
    "addr_info": []
  },
  {
    "ifindex": 3,
    "ifname": "wlp2s0",
    "flags": [
      "BROADCAST",
      "MULTICAST",
      "UP",
      "LOWER_UP"
    ],
    "mtu": 1500,
    "qdisc": "mq",
    "operstate": "UP",
    "group": "default",
    "txqlen": 1000,
    "link_type": "ether",
    "address": "00:f4:8d:7f:f0:81",
    "broadcast": "ff:ff:ff:ff:ff:ff",
    "addr_info": [
      {
        "family": "inet",
        "local": "192.168.1.12",
        "prefixlen": 24,
        "broadcast": "192.168.1.255",
        "scope": "global",
        "dynamic": true,
        "noprefixroute": true,
        "label": "wlp2s0",
        "valid_life_time": 81912,
        "preferred_life_time": 81912
      },
      {
        "family": "inet6",
        "local": "fe80::c5cb:a009:c6cb:14fc",
        "prefixlen": 64,
        "scope": "link",
        "noprefixroute": true,
        "valid_life_time": 4294967295,
        "preferred_life_time": 4294967295
      }
    ]
  },
  {
    "ifindex": 4,
    "ifname": "virbr0",
    "flags": [
      "NO-CARRIER",
      "BROADCAST",
      "MULTICAST",
      "UP"
    ],
    "mtu": 1500,
    "qdisc": "noqueue",
    "operstate": "DOWN",
    "group": "default",
    "txqlen": 1000,
    "link_type": "ether",
    "address": "52:54:00:43:4f:24",
    "broadcast": "ff:ff:ff:ff:ff:ff",
    "addr_info": [
      {
        "family": "inet",
        "local": "192.168.122.1",
        "prefixlen": 24,
        "broadcast": "192.168.122.255",
        "scope": "global",
        "label": "virbr0",
        "valid_life_time": 4294967295,
        "preferred_life_time": 4294967295
      }
    ]
  },
  {
    "ifindex": 5,
    "ifname": "virbr0-nic",
    "flags": [
      "BROADCAST",
      "MULTICAST"
    ],
    "mtu": 1500,
    "qdisc": "fq_codel",
    "master": "virbr0",
    "operstate": "DOWN",
    "group": "default",
    "txqlen": 1000,
    "link_type": "ether",
    "address": "52:54:00:43:4f:24",
    "broadcast": "ff:ff:ff:ff:ff:ff",
    "addr_info": []
  },
  {
    "ifindex": 6,
    "ifname": "docker0",
    "flags": [
      "NO-CARRIER",
      "BROADCAST",
      "MULTICAST",
      "UP"
    ],
    "mtu": 1500,
    "qdisc": "noqueue",
    "operstate": "DOWN",
    "group": "default",
    "link_type": "ether",
    "address": "02:42:7b:4d:6d:1e",
    "broadcast": "ff:ff:ff:ff:ff:ff",
    "addr_info": [
      {
        "family": "inet",
        "local": "172.17.0.1",
        "prefixlen": 16,
        "broadcast": "172.17.255.255",
        "scope": "global",
        "label": "docker0",
        "valid_life_time": 4294967295,
        "preferred_life_time": 4294967295
      }
    ]
  }
]`

func Test1(t *testing.T) {
	u, _ := model.NewSSHUserByPassphraseWithStringLine(`panda, 123456, 192.168.122.10:22`)
	term, err := m_terminal.GetSSHClientByPassphrase(u)
	if err != nil {
		t.Fatal(err)
	}
	m_terminal.GetRemoteHostInfo(term)
	content := term.GetContent()
	info, ok := content.GetHostInfo()
	if ok {
		u := info.User
		fmt.Printf("%d\n", u.Uid)

	}
}

func Test2(t *testing.T) {
	extractName, _ := regexp.Compile(`\d+\((\w+)\)`)
	str := `1000(pengda),4(adm),24(cdrom),27(sudo),30(dip),46(plugdev),121(lpadmin),131(lxd),132(sambashare),135(libvirt),998(docker)`
	for _, i := range strings.Split(str, ",") {
		t := extractName.FindStringSubmatch(i)
		fmt.Println(t[1])
	}
}
