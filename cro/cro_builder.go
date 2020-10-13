package cro

import (
	"io"
	"time"
)

type (
	taskBuilder struct {
		// 主机选择器
		hostsF string
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

func (c *taskBuilder) Timeout(t time.Duration) *taskBuilder {
	c.timeout = t
	return c
}

func (c *taskBuilder) Format(format string) *taskBuilder {
	c.format = format
	return c
}

func (c *taskBuilder) ListenP(port int) *taskBuilder {
	c.listenP = uint16(port)
	return c
}

func (c *taskBuilder) Out(o io.Writer) *taskBuilder {
	c.out = o
	return c
}

func (c *taskBuilder) Filter() *taskBuilder {
	return c
}
