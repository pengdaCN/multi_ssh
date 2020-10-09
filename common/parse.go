package common

import (
	"regexp"
	"strings"
)

const (
	chg    = '/'
	border = `'"`
)

var (
	firstChar      = regexp.MustCompile(`(?:\s+)?.`)
	endWithSpace   = regexp.MustCompile(`.*?(?:\s)`)
	endWithSingleB = regexp.MustCompile(`(?:').*?(?:')`)
	endWithDoubleB = regexp.MustCompile(`(?:").*?(?:")`)
)

func ReadStr(src string) (str string, stream string) {
	var (
		b   string
		rst strings.Builder
	)
	f := firstChar.FindString(src)
	{
		f = string(f[len(f)-1])
	}
	src = src[len(f):]
	if strings.Contains(border, f) {
		b = f
	}
	for {
		var s string
		switch b {
		case `'`:
			s = endWithSingleB.FindString(src)
		case `"`:
			s = endWithDoubleB.FindString(src)
		default:
			s = endWithSpace.FindString(src)
		}
		if b == "" {
			src = src[len(s):]
		} else {
			src = src[len(s)+1:]
		}
		// 使用空格分割是没有转义功能
		if b == "" {
			str = s
			break
		}
		var _f bool
		// 存在转义字符
		for {
			i := strings.IndexRune(s, chg)
			if i == -1 {
				break
			}
			rst.WriteString(s[:i])
			if i == len(s)-1 {
				s = ""
				_f = true
				rst.WriteString(b)
				break
			}
			switch s[i+1] {
			case 'n':
				rst.WriteRune('\n')
			case 'r':
				rst.WriteRune('\r')
			case 'v':
				rst.WriteRune('\v')
			case '\\':
				rst.WriteRune('\\')
			case 't':
				rst.WriteRune('\t')
			case 'b':
				rst.WriteRune('\b')
			case 'f':
				rst.WriteRune('\f')
			default:
				panic("ERROR Bad escape character")
			}
			s = s[i+1:]
		}
		rst.WriteString(s)
		if !_f {
			str = rst.String()
			str = strings.Trim(str, b)
			break
		}
	}
	return
}
