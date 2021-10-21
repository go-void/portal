package cli

type Context struct {
	App     *App
	Command *Command
	Args    []string
}

func (c *Context) HasArgs() bool {
	return len(c.Args) != 0
}
