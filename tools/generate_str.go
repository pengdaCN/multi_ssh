package tools

import (
	cr "crypto/rand"
	"encoding/binary"
	"math/rand"
	"strings"
)

const (
	word      = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789+="
	bitLength = 6
	bitMask   = 1<<bitLength - 1
)

var (
	Ra rand.Source
)

func init() {
	var r [4]byte
	// 获取随机数种子
	_, _ = cr.Reader.Read(r[:])
	seed, _ := binary.Varint(r[:])
	Ra = rand.NewSource(seed)
}

// param n 生成字符个数
// 生成随机的个数的base64字符
func GenerateRandomStr(n int) string {
	var s strings.Builder
	for i, cache := 0, Ra.Int63(); i < n; i++ {
		if cache == 0 {
			cache = Ra.Int63()
		}
		idx := cache & bitMask
		s.WriteByte(word[idx])
		cache >>= bitLength
	}
	return s.String()
}
