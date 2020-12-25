package cro

type playbookTBuilder struct {
	baseBuilder *baseTaskBuilder
	path        string
	text        string
	vars        string
}

func (c *playbookTBuilder) SetVars(vars string) *playbookTBuilder {
	c.vars = vars
	return c
}

func (c *playbookTBuilder) Path(path string) *playbookTBuilder {
	c.path = path
	return c
}

func (c *playbookTBuilder) Text(text string) *playbookTBuilder {
	c.text = text
	return c
}

func (c *playbookTBuilder) Builder() *playbookTRunEnv {
	return &playbookTRunEnv{
		conf: c,
		b:    c.baseBuilder.builder(),
	}
}
