package common

import (
	"github.com/bluele/gcache"
	"regexp"
)

// 正则表达式缓存
var ReCache = gcache.New(50).LRU().Build()

func GetRe(re string) *regexp.Regexp {
	if v, found := ReCache.Get(re); found != nil {
		reb := regexp.MustCompile(re)
		_ = ReCache.Set(re, reb)
		return reb
	} else {
		return v.(*regexp.Regexp)
	}
}
