package cro

type shellTBuilder struct {
	baseBuilder *baseTaskBuilder
	sudo        bool
	cmds        string
}

func (c *shellTBuilder) Sudo(s bool) *shellTBuilder {
	c.sudo = s
	return c
}

func (c *shellTBuilder) Cmds(args string) *shellTBuilder {
	c.cmds = args
	return c
}

func (c *shellTBuilder) Builder() *shellRunEnv {
	return &shellRunEnv{
		conf: c,
		b:    c.baseBuilder.builder(),
	}
}
