package lexer

import (
	"io"
	"io/ioutil"
	"os"
	"path"
)

type (
	parseStat int
	parser    struct {
		chunk     string
		curLine   int
		parseStat parseStat
		filename  string
	}
	ParserOption struct {
		Filename string
	}
)

const (
	started parseStat = iota
	overed
)

func NewParser(file string) *parser {
	f, err := os.Open(file)
	if err != nil {
		panic(err)
	}
	filename := path.Base(file)
	return newParser(f, &ParserOption{
		Filename: filename,
	})
}

func newParser(reader io.Reader, option *ParserOption) *parser {
	c, err := ioutil.ReadAll(reader)
	if err != nil {
		panic(err)
	}
	var (
		filename string
	)
	if option != nil && option.Filename != "" {
		filename = option.Filename
	} else {
		filename = "built_main"
	}
	return &parser{
		chunk:    string(c),
		filename: filename,
		curLine:  1,
	}
}
