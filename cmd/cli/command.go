package cli

type Command struct {
	Name   string
	Usage  string
	Action ActionFunc
}

func (c *Command) Run(ctx *Context) error {
	return c.Action(ctx)
}
