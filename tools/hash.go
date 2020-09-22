package tools

import (
	"crypto/md5"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"os"
)

const Gb = 1 << (10 * 3)

func HashFile(fileName string) (hStr string, err error) {
	fStat, err := os.Stat(fileName)
	if err != nil {
		return "", err
	}
	if fStat.Size() > 10*Gb {
		j := []byte(fmt.Sprintf("%s %d %s %s", fStat.Name(), fStat.Size(), fStat.Mode(), fStat.ModTime()))

		h := md5.Sum(j)
		return base64.StdEncoding.EncodeToString(h[:]), nil
	}
	b, err := ioutil.ReadFile(fileName)
	if err != nil {
		return "", err
	}
	h := md5.Sum(b)
	return base64.StdEncoding.EncodeToString(h[:]), nil
}
