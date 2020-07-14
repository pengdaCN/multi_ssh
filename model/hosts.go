package model

type RemoteHostInfo struct {
	UserName   string
	Passphrase string
	Host       string
	Alias      string
}

//func LoadHost(str string) (*RemoteHostInfo, error) {
//	s := strings.TrimSpace(str)
//	piece := separate.Split(str, -1)
//}
