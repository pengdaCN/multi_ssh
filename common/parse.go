package common

import (
	"regexp"
	"strings"
)

const (
	chg    rune = '\\'
	border      = `'"`
)

var (
	firstChar = regexp.MustCompile(`(?:\s+)?.`)
	//endWithSpace   = regexp.MustCompile(`[^\s]*`)
	endWithSingleB = regexp.MustCompile(`.*?(?:')`)
	endWithDoubleB = regexp.MustCompile(`.*?(?:")`)
	firstWord      = regexp.MustCompile(`^\s*[a-zA-Z_](?:\w)*`)
)

func ReadStr(src string) (str string, stream string) {
	var (
		b     string
		bchar byte
		rst   strings.Builder
	)
	f := firstChar.FindString(src)
	l := len(f)
	if l == 0 {
		return "", ""
	}
	{
		f = string(f[len(f)-1])
	}
	src = src[l:]
	if strings.Contains(border, f) {
		b = f
		bchar = b[0]
	} else {
		panic("ERROR unknown border char")
	}
	var (
		useSearchWord *regexp.Regexp
	)
	switch b {
	case `'`:
		useSearchWord = endWithSingleB
	case `"`:
		useSearchWord = endWithDoubleB
	}
	var s string
TAKE:
	s = useSearchWord.FindString(src)
	if s == "" {
		panic("ERROR")
	}
	src = src[len(s):]
	for {
		i := strings.IndexRune(s, chg)
		if i == -1 {
			rst.WriteString(s[:len(s)-1])
			break
		}
		rst.WriteString(s[:i])
		switch s[i+1] {
		case bchar:
			rst.WriteByte(bchar)
			if i == len(s)-1-1 {
				goto TAKE
			}
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
		s = s[i+2:]
	}
	return rst.String(), src
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
