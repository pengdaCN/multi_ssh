package main

import (
	"multi_ssh/m_terminal"
	"multi_ssh/model"
	"testing"
)

func Test1(t *testing.T) {
	u := model.SSHUserByPassphrase{
		UserName: "serv",
		Password: "123456",
		RemoteHost: "192.168.101.16:22",
	}
	client, _ := m_terminal.GetSSHClientByPassphrase(&u)
	_ = client.SftpUpdate("/home/pengda/myCode/go/multi_ssh/main.go", "/home/serv")
	//println(err.Error())
}
