package tools

import (
	"reflect"
	"unsafe"
)

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
