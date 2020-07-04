package main

import (
	"log"
	"multi_ssh/m_terminal"
	"multi_ssh/model"
	"sync"
	"time"
)

type candy struct {
	u       model.SHHUser
	cmd     string
	timeout time.Duration
}

func exec(out chan<- pack, in ...*candy) {
	defer close(out)
	var w sync.WaitGroup
	for _, i := range in {
		w.Add(1)
		go func(c *candy) {
			defer w.Done()
			client, err := m_terminal.GetSSHClientByPassphrase(c.u)
			defer func() {
				_ = client.Close()
			}()
			if err != nil {
				log.Println("链接失败", c.u.User(), err.Error())
				return
			}
			if b, err := client.Run(true, c.cmd); err == nil {
				out <- pack{
					msg:      b,
					time:     time.Now(),
					userInfo: c.u,
				}
			} else {
				out <- pack{
					msg:      []byte(err.Error()),
					time:     time.Now(),
					userInfo: c.u,
				}
			}
		}(i)
	}
	w.Wait()
}

func exec2(out chan<- pack, in ...*candy) {
	defer close(out)
	var w sync.WaitGroup
	for _, i := range in {
		w.Add(1)
		go func(c *candy) {
			defer w.Done()
			client, err := m_terminal.GetSSHClientByPassphrase(c.u)
			defer func() {
				_ = client.Close()
			}()
			if err != nil {
				log.Println("链接失败", c.u.User(), err.Error())
				return
			}
			if b, err := client.Run2(c.cmd); err == nil {
				out <- pack{
					msg:      b,
					time:     time.Now(),
					userInfo: c.u,
				}
			} else {
				out <- pack{
					msg:      []byte(err.Error()),
					time:     time.Now(),
					userInfo: c.u,
				}
			}
		}(i)
	}
	w.Wait()
}
