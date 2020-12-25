package cro

type copyTBuilder struct {
	baseBuilder *baseTaskBuilder
	src         []string
	dst         string
	sudo        bool
	exists      bool
}

func (c *copyTBuilder) Sudo(sudo bool) *copyTBuilder {
	c.sudo = sudo
	return c
}

func (c *copyTBuilder) Exists(exists bool) *copyTBuilder {
	c.exists = exists
	return c
}

func (c *copyTBuilder) Src(src []string) *copyTBuilder {
	c.src = src
	return c
}

func (c *copyTBuilder) Dst(dst string) *copyTBuilder {
	c.dst = dst
	return c
}

func (c *copyTBuilder) Builder() *copyTRunEnv {
	return &copyTRunEnv{
		conf: c,
		b:    c.baseBuilder.builder(),
	}
}
