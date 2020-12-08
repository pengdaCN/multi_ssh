package common

import (
	"fmt"
	"regexp"
	"strings"
)

const (
	chg rune = '\\'
)

var (
	ReadBetween = readRuneBetween
)

var (
	firstChar = regexp.MustCompile(`(?:\s+)?.`)
	firstWord = regexp.MustCompile(`^\s*[a-zA-Z_](?:\w)*`)
)

func ReadStr(src string) (str string, stream string) {
	var (
		single = [2]rune{'\'', '\''}
		double = [2]rune{'"', '"'}
	)
	if str, stream = readRuneBetween(src, single); str != "" {
		return
	}
	str, stream = readRuneBetween(src, double)
	return
}

func ReadNotSpaceChar(src string) (str, stream string) {
	word := firstWord.FindString(src)
	str = strings.TrimSpace(word)
	stream = src[len(word):]
	return
}

func ReadWord(src string) (word string, stream string) {
	word = firstWord.FindString(src)
	if word == "" {
		panic("ERROR")
	}
	stream = src[len(word):]
	word = strings.TrimSpace(word)
	return
}

func expandCharts(sb *strings.Builder, str string) {
	for {
		i := strings.IndexRune(str, chg)
		if i < 0 {
			sb.WriteString(str)
			return
		}
		if i < len(str)-1 {
			sb.WriteString(str[:i])
			switch str[i+1] {
			case 'n':
				sb.WriteRune('\n')
			case 'r':
				sb.WriteRune('\r')
			case 'v':
				sb.WriteRune('\v')
			case '\\':
				sb.WriteRune('\\')
			case 't':
				sb.WriteRune('\t')
			case 'b':
				sb.WriteRune('\b')
			case 'f':
				sb.WriteRune('\f')
			default:
				panic("ERROR Bad escape character")
			}
			str = str[i+1:]
		} else {
			panic("ERROR Bad escape character")
		}
	}

}

// 起至和结束符号不能为/ 与空白字符
func readRuneBetween(src string, symbol [2]rune) (rst, stream string) {
	firstCharts := firstChar.FindString(src)
	if firstCharts == "" {
		return "", src
	}
	if !strings.HasSuffix(firstCharts, string(symbol[0])) {
		return "", src
	}
	src = src[len(firstCharts):]
	var sb strings.Builder
	endS := fmt.Sprintf(`\%s`, string(symbol[1]))
WALK:
	i := strings.IndexRune(src, symbol[1])
	if i < 0 {
		panic("ERROR no normal end")
	}
	if strings.HasSuffix(src[:i+1], endS) {
		expandCharts(&sb, src[:i-len(endS)+1])
		sb.WriteRune(symbol[1])
		src = src[i+1:]
		goto WALK
	}
	sb.WriteString(src[:i])
	src = src[i:]
	rst = sb.String()
	return
}
