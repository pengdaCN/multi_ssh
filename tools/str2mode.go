package tools

import (
	"errors"
	"os"
)

/*
perm 格式 rwxrwxrwx
*/

var str2perm map[rune]uint32

func init() {
	str2perm = map[rune]uint32{}
	str2perm['x'] = 0
	str2perm['w'] = 1
	str2perm['r'] = 2
}
func String2FileMode(perm string) (os.FileMode, error) {
	if len(perm) != 9 {
		return 0, errors.New("perm 格式不对")
	}
	var permission os.FileMode
	for i, p := range perm {
		n, ok := str2perm[p]
		if !ok {
			continue
		}
		offset := n + uint32(i/3)*3
		if i%3 != 0 {
			//	不允许权限重复的设置
			t := permission & (1 << offset)
			if t != 0 {
				return 0, errors.New("ERROR 重复设置权限")
			}
		}
		permission |= 1 << offset
	}
	return permission, nil
}
