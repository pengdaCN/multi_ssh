package cro

import (
	"io"
	"io/ioutil"
	"multi_ssh/model"
	"multi_ssh/tools"
	"time"
)

type (
	hostsFilter func(user model.SHHUser) bool
	croBuilder  struct {
		// 主机选择器
		hostsF hostsFilter
		// 设置执行任务最长时间
		timeout time.Duration
		// 原始的主机信息
		rawHostsInfo string
		// 输出信息格式
		format string
		// 执行时监听端口
		listenP uint16
		// 任务信息输出位置
		out io.Writer
	}
)

func ReadF(filename string) *croBuilder {
	f, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	return &croBuilder{
		rawHostsInfo: tools.ByteSlice2String(f),
	}
}

func ReadS(hostsInfo string) *croBuilder {
	return &croBuilder{
		rawHostsInfo: hostsInfo,
	}
}

func (c *croBuilder) Timeout(t time.Duration) *croBuilder {
	c.timeout = t
	return c
}

func (c *croBuilder) Format(format string) *croBuilder {
	c.format = format
	return c
}

func (c *croBuilder) ListenP(port int) *croBuilder {
	c.listenP = uint16(port)
	return c
}

func (c *croBuilder) Out(o io.Writer) *croBuilder {
	c.out = o
	return c
}

func (c *croBuilder) Filter() *croBuilder {
	return c
}
