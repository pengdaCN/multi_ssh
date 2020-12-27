package cro

import (
	"multi_ssh/extra_mod/pool"
	"multi_ssh/model"
	"net"
	"strconv"
	"strings"
)

const (
	unknown op = iota
	equal
	noEqual
	borderSymbol = `'"`
	chgWord      = '\\'
)

var (
	opM = map[op]string{
		equal:   "==",
		noEqual: "!=",
	}
	s2op  map[string]op
	opStr []string
	pM    = map[string]matchFunc{
		"IP":   matchIP,
		"USER": matchUser,
		"PORT": matchPort,
	}
)

func init() {
	opStr = make([]string, 0, len(opM))
	for _, v := range opM {
		opStr = append(opStr, v)
	}
	s2op = make(map[string]op, len(opM))
	for k, v := range opM {
		s2op[v] = k
	}
}

type (
	op    uint8
	token struct {
		rightOP string
		OP      op
		liftOP  string
	}
	matchFunc func(model.SHHUser, *token) bool
	piece     struct {
		fn matchFunc
		t  *token
	}
)

func filters(src []model.SHHUser, str string) []model.SHHUser {
	fi := parseFilter(str)
	rst := make([]model.SHHUser, 0)
	for _, v := range src {
		if filter(v, fi) {
			rst = append(rst, v)
		}
	}
	return rst
}

func filter(src model.SHHUser, fi []piece) bool {
	for _, v := range fi {
		if !v.fn(src, v.t) {
			return false
		}
	}
	return true
}

const (
	rangeKeyword = '-'
	manyKeyword  = ','
	netKeyword   = '/'
)

type (
	ipHandle func(net.IP, string) bool
)

var (
	mIPm = map[rune]ipHandle{
		-1:           matchIPSingle,
		rangeKeyword: matchIPManyRange,
		manyKeyword:  matchIPManyRange,
		netKeyword:   matchIPNet,
	}
)

func matchPort(h model.SHHUser, t *token) bool {
	var (
		p1 bool
		p2 bool
	)
	_, port, err := net.SplitHostPort(h.Host())
	if err != nil {
		panic("解析ip:port错误")
	}
	if port == t.liftOP {
		p1 = true
	}
	p2 = t.OP == equal
	if p2 {
		return p1
	}
	return !p1
}

func matchUser(h model.SHHUser, t *token) bool {
	var (
		p1 bool
		p2 bool
	)
	if t.liftOP == h.User() {
		p1 = true
	}
	p2 = t.OP == equal
	if p2 {
		return p1
	}
	return !p1
}

func matchIP(h model.SHHUser, t *token) bool {
	var (
		p1   bool
		p2   bool
		mode rune
	)
	ipStr, _, err := net.SplitHostPort(h.Host())
	if err != nil {
		panic("解析ip:port错误")
	}
	ip := net.ParseIP(ipStr)
	switch {
	case strings.ContainsRune(t.liftOP, netKeyword):
		mode = netKeyword
	case strings.ContainsRune(t.liftOP, rangeKeyword):
		mode = rangeKeyword
	case strings.ContainsRune(t.liftOP, manyKeyword):
		mode = manyKeyword
	default:
		mode = -1
	}
	if fn, ok := mIPm[mode]; ok {
		p1 = fn(ip, t.liftOP)
	}
	p2 = t.OP == equal
	if p2 {
		return p1
	}
	return !p1
}

func matchIPSingle(ip net.IP, filstr string) bool {
	var (
		tIP net.IP
	)
	if v, ok := pool.Share.Load(filstr); ok {
		tIP = v.(net.IP)
	} else {
		tIP = net.ParseIP(filstr)
		pool.Share.Store(filstr, tIP)
	}
	return tIP.Equal(ip)
}

type (
	pair struct {
		start uint8
		end   uint8
	}
	pairs []pair

	ipRange struct {
		ip      [4]byte
		idx     uint8
		handles [4]pairs
	}
)

func newIpRange(str string) *ipRange {
	ipPairs := strings.Split(str, ".")
	if len(ipPairs) != 4 {
		panic("错误的ipv4格式")
	}
	var (
		parts   [4]byte
		handles [4]pairs
		idx     uint8
	)
	for i, v := range ipPairs {
		if strings.ContainsRune(v, rangeKeyword) || strings.ContainsRune(v, manyKeyword) {
			p := parseManyRange(v)
			if p == nil {
				panic("错误many range 格式")
			}
			handles[i] = p
			idx |= 1 << i
			continue
		}
		t, err := strconv.ParseUint(v, 10, 8)
		if err != nil {
			panic("错误的ipv4 格式，不能解析数字")
		}
		part := byte(t)
		parts[i] = part
	}
	return &ipRange{
		ip:      parts,
		handles: handles,
		idx:     idx,
	}
}

func parseManyRange(str string) pairs {
	p := make(pairs, 0)
	ps := strings.Split(str, string(manyKeyword))
	for _, v := range ps {
		if strings.ContainsRune(v, rangeKeyword) {
			part := strings.Split(v, string(rangeKeyword))
			if len(part) != 2 {
				return nil
			}
			s, err := strconv.ParseUint(part[0], 10, 8)
			if err != nil {
				return nil
			}
			start := byte(s)
			e, err := strconv.ParseUint(part[1], 10, 8)
			if err != nil {
				return nil
			}
			end := byte(e)
			if start > end {
				return nil
			}
			p = append(p, pair{start: start, end: end})
			continue
		}
		i, err := strconv.ParseUint(v, 10, 8)
		if err != nil {
			return nil
		}
		_i := byte(i)
		p = append(p, pair{start: _i, end: _i})
	}
	return p
}

func (p *ipRange) contain(ip net.IP) bool {
	parts := []byte(ip.To4())
	for i, v := range parts[:4] {
		e := p.idx & (1 << i)
		if e != 0 {
			if !p.handles[i].contain(v) {
				return false
			}
			continue
		}
		if p.ip[i] != v {
			return false
		}
	}
	return true
}

func (p *pairs) contain(part byte) bool {
	for _, v := range *p {
		if v.start <= part && v.end >= part {
			return true
		}
	}
	return false
}

func matchIPManyRange(ip net.IP, filstr string) bool {
	var (
		iprange *ipRange
	)
	if v, ok := pool.Share.Load(filstr); ok {
		iprange = v.(*ipRange)
	} else {
		iprange = newIpRange(filstr)
		pool.Share.Store(filstr, iprange)
	}
	return iprange.contain(ip)
}

func matchIPNet(ip net.IP, filstr string) bool {
	var (
		ipNET *net.IPNet
	)
	if v, ok := pool.Share.Load(filstr); ok {
		ipNET = v.(*net.IPNet)
	} else {
		_, _ipNET, err := net.ParseCIDR(filstr)
		if err != nil {
			panic("解析ip 网段")
		}
		ipNET = _ipNET
		pool.Share.Store(filstr, ipNET)
	}
	return ipNET.Contains(ip)
}

func matchExtraInfo(h model.SHHUser, t *token) bool {
	info := h.Extra()
	var (
		p1 bool
		p2 bool
	)
	p2 = t.OP == equal
	if v, ok := info[t.rightOP]; ok {
		if v == t.liftOP {
			p1 = true
		} else {
			p1 = false
		}
	}
	if p2 {
		return p1
	}
	return !p1
}

const (
	step0 = iota
	step1
	step2
	step3
	step4
	step5
)

func parseFilter(str string) []piece {
	var (
		state      = step0
		rst        = make([]piece, 0)
		rightPiece strings.Builder
		_op        string
		p          op
		leftPiece  strings.Builder
		border     rune
	)
	for _, v := range str {
		switch state {
		case step0:
			if v == ' ' || v == '\n' || v == '\r' || v == '\t' {
				continue
			}
			if v >= 'A' && v <= 'Z' || v >= 'a' && v <= 'z' {
				rightPiece.WriteRune(v)
				state = step1
				continue
			}
			panic("解析错误")
		case step1:
			if v >= 'A' && v <= 'Z' || v >= 'a' && v <= 'z' || v >= '0' && v <= '9' {
				rightPiece.WriteRune(v)
				continue
			}
			_op += string(v)
			state = step2
		case step2:
			_op += string(v)
			if _p, ok := matchOP(_op); ok {
				if _p != unknown {
					p = _p
					state = step3
					continue
				}
				continue
			}
			panic("错误的op")
		case step3:
			if strings.ContainsRune(borderSymbol, v) {
				border = v
				state = step4
				continue
			} else {
				border = ' '
			}
		case step4:
			if v == chgWord {
				state = step5
				continue
			}
			if v == border {
				state = step0
				r := rightPiece.String()
				if r == "" {
					panic("错误的右值")
				}
				l := leftPiece.String()
				{
					t := &token{
						rightOP: r,
						OP:      p,
						liftOP:  l,
					}
					pi := piece{
						t: t,
					}
					if p, ok := pM[r]; ok {
						pi.fn = p
					} else {
						pi.fn = matchExtraInfo
					}
					rst = append(rst, pi)
				}
				{
					rightPiece = strings.Builder{}
					_op = ""
					p = unknown
					leftPiece = strings.Builder{}
					border = 0
				}
				continue
			}
			leftPiece.WriteRune(v)
		case step5:
			leftPiece.WriteRune(v)
			state = step4
		}
	}
	if state != step0 {
		panic("错误的结束filter 格式匹配字符串")
	}
	return rst
}

func matchOP(str string) (op, bool) {
	if p, ok := s2op[str]; ok {
		return p, ok
	}
	for _, v := range opStr {
		if strings.HasPrefix(v, str) {
			return unknown, true
		}
	}
	return unknown, false
}
