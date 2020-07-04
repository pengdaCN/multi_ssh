package model

import "golang.org/x/crypto/ssh"

type SHHUser interface {
	Host() string
	User() string
	Auth() []ssh.AuthMethod
}

type SSHUserByPassphrase struct {
	RemoteHost string
	UserName   string
	Password   string
}

func (s *SSHUserByPassphrase) Host() string {
	return s.RemoteHost
}

func (s *SSHUserByPassphrase) User() string {
	return s.UserName
}

func (s *SSHUserByPassphrase) Auth() []ssh.AuthMethod {
	return []ssh.AuthMethod{ssh.Password(s.Password)}
}
