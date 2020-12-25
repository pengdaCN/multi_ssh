package cro

type scriptTBuilder struct {
	baseBuilder *baseTaskBuilder
	sudo        bool
	args        string
	path        string
	text        string
}

func (c *scriptTBuilder) Sudo(s bool) *scriptTBuilder {
	c.sudo = s
	return c
}

func (c *scriptTBuilder) Args(args string) *scriptTBuilder {
	c.args = args
	return c
}

func (c *scriptTBuilder) Path(path string) *scriptTBuilder {
	c.path = path
	return c
}

func (c *scriptTBuilder) Text(text string) *scriptTBuilder {
	c.text = text
	return c
}

func (c *scriptTBuilder) Builder() *ScriptEunEnv {
	return &ScriptEunEnv{
		conf: c,
		b:    c.baseBuilder.builder(),
	}
}
