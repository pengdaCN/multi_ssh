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
	if str, stream = readRuneBetween(src, single, true); str != "" {
		return
	}
	str, stream = readRuneBetween(src, double, true)
	return
}

func ReadNotSpaceChar(src string) (str, stream string) {
	word := firstChar.FindString(src)
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
func readRuneBetween(src string, symbol [2]rune, expand bool) (rst, stream string) {
	stream = src
	firstCharts := firstChar.FindString(stream)
	if firstCharts == "" {
		return "", stream
	}
	if !strings.HasSuffix(firstCharts, string(symbol[0])) {
		return "", stream
	}
	stream = stream[len(firstCharts):]
	var sb strings.Builder
	endS := fmt.Sprintf(`\%s`, string(symbol[1]))
WALK:
	i := strings.IndexRune(stream, symbol[1])
	if i < 0 {
		panic("ERROR no normal end")
	}
	if strings.HasSuffix(stream[:i+1], endS) {
		if expand {
			expandCharts(&sb, stream[:i-len(endS)+1])
		} else {
			sb.WriteString(stream[:i-len(endS)+1])
		}
		sb.WriteRune(symbol[1])
		stream = stream[i+1:]
		goto WALK
	}
	if expand {
		expandCharts(&sb, stream[:i])
	} else {
		sb.WriteString(stream[:i])
	}
	stream = stream[i+1:]
	rst = sb.String()
	return
}
