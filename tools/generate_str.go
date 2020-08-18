package tools

import (
	"math/rand"
	"strings"
	"time"
)

const (
	word      = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789+="
	bitLength = 6
	bitMask   = 1<<bitLength - 1
)

var (
	ra = rand.NewSource(time.Now().UnixNano() + 5052)
)

// param n 生成字符个数
// 生成随机的个数的base64字符
func GenerateRandomStr(n int) string {
	var s strings.Builder
	for i, cache := 0, ra.Int63(); i < n; i++ {
		if cache == 0 {
			cache = ra.Int63()
		}
		idx := cache & bitMask
		s.WriteByte(word[idx])
		cache >>= bitLength
	}
	return s.String()
}
