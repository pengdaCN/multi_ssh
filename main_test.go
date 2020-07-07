package main

import (
	"fmt"
	"multi_ssh/m_terminal"
	"multi_ssh/model"
	"testing"
)

func Test1(t *testing.T) {
	u := model.SSHUserByPassphrase{
		UserName:   "panda",
		Password:   "123456",
		RemoteHost: "192.168.122.10:22",
	}
	client, err := m_terminal.GetSSHClientByPassphrase(&u)
	if err != nil {
		t.Fatal(err.Error())
	}
	bs, err := client.Run(true, "sudo whoami")
	if err != nil {
		t.Fatal(err.Error())
	}
	fmt.Println(string(bs))
}