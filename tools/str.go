package tools

import (
	"math/rand"
	"reflect"
	"time"
	"unsafe"
)

func init() {
	rand.Seed(time.Now().UnixNano())
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
