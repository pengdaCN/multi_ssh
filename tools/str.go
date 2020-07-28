package tools

import (
	"math/rand"
	"reflect"
	"strings"
	"time"
	"unsafe"
)

var specialStr map[rune]rune

func init() {
	rand.Seed(time.Now().UnixNano())
	specialStr = make(map[rune]rune)
	specialStr['n'] = 10
	specialStr['t'] = 9
}

func ByteSlice2String(b []byte) string {
	slice := (*reflect.SliceHeader)(unsafe.Pointer(&b))
	str := reflect.StringHeader{
		Data: slice.Data,
		Len:  slice.Len,
	}
	return *(*string)(unsafe.Pointer(&str))
}

func String2ByteSlice(s string) []byte {
	str := (*reflect.StringHeader)(unsafe.Pointer(&s))
	slice := reflect.SliceHeader{
		Data: str.Data,
		Len:  str.Len,
		Cap:  str.Len,
	}
	return *(*[]byte)(unsafe.Pointer(&slice))
}

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func RandStringBytes(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

// 转换特殊字符，如\n \t
// 目前仅支持\n \t
func SpecialStrTransform(str string) string {
	var special uint8
	var rst strings.Builder
	for _, v := range str {
		if v == '\\' {
			special ^= 1
			continue
		}
		if special == 1 {
			if val, ok := specialStr[v]; ok {
				rst.WriteRune(val)
				special &= 0
				continue
			}
		}
		rst.WriteRune(v)
	}
	return rst.String()
}
