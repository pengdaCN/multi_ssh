package main

import (
	"fmt"
	"multi_ssh/tools"
	"testing"
)

func Test1(t *testing.T) {
	s := "r-xrwxr-x"
	mode, err := tools.String2FileMode(s)
	fmt.Printf("%b", mode)
	if err != nil {
		fmt.Println(err.Error())
	}
}
