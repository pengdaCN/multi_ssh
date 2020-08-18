package model

import "regexp"

var (
	separate, _ = regexp.Compile(`\s*,\s*`)
	ignoreLine, _ = regexp.Compile(`^\s*#`)
	spaceLine, _ = regexp.Compile(`^\s+$`)
	extraPiece, _ = regexp.Compile("(.*) +(`.*`)")
)

type RemoteHostInfo struct {
	UserName   string
	Passphrase string
	Host       string
	Alias      string
	Extra      string
}

func ParseLine(line string) *RemoteHostInfo {
	if ignoreLine.MatchString(line) || ignoreLine.MatchString(line) {
		return nil
	}
	sli := extraPiece.FindStringSubmatch(line)
	if len(sli) < 2 {
		return nil
	}
	r := new(RemoteHostInfo)
	if ! parseBase(r, sli[1]) {
		return nil
	}
	if len(sli) >= 3 {
		r.Extra = sli[2]
	}
	return r
}

func parseBase(r *RemoteHostInfo, str string) bool {
	arr := separate.Split(str, -1)
	if len(arr) != 3 {
		return false
	}
	r.UserName = arr[0]
	r.Passphrase = arr[1]
	r.Host = arr[2]
	return true
}
