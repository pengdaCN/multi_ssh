package cmd

import (
	"multi_ssh/model"
	"strings"
)

const (
	unknown op = iota
	equal
	noEqual
	borderSign = `'"`
	chgWord    = '\\'
)

var (
	opM = map[op]string{
		equal:   "==",
		noEqual: "!=",
	}
	s2op  map[string]op
	opStr []string
	pM    = map[string]matchFunc{
		"ip": matchIP,
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

func matchIP(h model.SHHUser, t *token) bool {

	return false
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
			if strings.ContainsRune(borderSign, v) {
				border = v
				state = step4
				continue
			}
			panic("错误的border符号")
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
