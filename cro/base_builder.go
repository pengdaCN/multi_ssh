package cro

import (
	"io"
	"time"
)

type baseTaskBuilder struct {
	// 主机选择器
	hostsF string
	// 设置执行任务最长时间
	timeout time.Duration
	// 原始的主机信息
	rawHostsInfo string
	// 输出信息格式
	format   string
	filerStr string
	// 任务信息输出位置
	out            io.Writer
	execMaxSeveral int
	preInfo        bool
}

func New() *baseTaskBuilder {
	return new(baseTaskBuilder)
}

func (c *baseTaskBuilder) Timeout(t time.Duration) *baseTaskBuilder {
	c.timeout = t
	return c
}

func (c *baseTaskBuilder) Format(format string) *baseTaskBuilder {
	c.format = format
	return c
}

func (c *baseTaskBuilder) Out(o io.Writer) *baseTaskBuilder {
	c.out = o
	return c
}

func (c *baseTaskBuilder) Filter(f string) *baseTaskBuilder {
	c.filerStr = f
	return c
}

func (c *baseTaskBuilder) Hosts(path string) *baseTaskBuilder {
	c.hostsF = path
	return c
}

func (c *baseTaskBuilder) Line(l string) *baseTaskBuilder {
	c.rawHostsInfo = l
	return c
}

func (c *baseTaskBuilder) RawHostsInfo(i string) *baseTaskBuilder {
	c.rawHostsInfo = i
	return c
}

func (c *baseTaskBuilder) PreInfo(p bool) *baseTaskBuilder {
	c.preInfo = p
	return c
}

func (c *baseTaskBuilder) SetMaxExecSeveral(limit int) *baseTaskBuilder {
	c.execMaxSeveral = limit
	return c
}

func (c *baseTaskBuilder) builder() *baseRunEnv {
	b, err := getBaseRunEnvFromBaseBuilder(c)
	if err != nil {
		panic(err)
	}
	return b
}

func (c *baseTaskBuilder) NewShellBuilder() *shellTBuilder {
	return &shellTBuilder{
		baseBuilder: c,
	}
}

func (c *baseTaskBuilder) NewScriptBuilder() *scriptTBuilder {
	return &scriptTBuilder{
		baseBuilder: c,
	}
}

func (c *baseTaskBuilder) NewPlaybookBuilder() *playbookTBuilder {
	return &playbookTBuilder{
		baseBuilder: c,
	}
}

func (c *baseTaskBuilder) NewCopyBuilder() *copyTBuilder {
	return &copyTBuilder{
		baseBuilder: c,
	}
}

func (c *baseTaskBuilder) NewPingRunEnv() *pingTRunEnv {
	return &pingTRunEnv{
		b: c.builder(),
	}
}
