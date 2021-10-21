package cli

type App struct {
	Name     string
	Usage    string
	Commands []*Command
	Action   ActionFunc
}

type ActionFunc func(*Context) error

func (a *App) Run(args []string) error {
	set := a.extractArgs(args)

	ctx := &Context{
		App:     a,
		Command: nil,
		Args:    set,
	}

	if ctx.HasArgs() {
		cmd := a.FindCommand(ctx.Args[0])
		if cmd != nil {
			return cmd.Run(ctx)
		}
	}

	return a.Action(ctx)
}

func (a *App) FindCommand(name string) *Command {
	for _, c := range a.Commands {
		if c.Name == name {
			return c
		}
	}
	return nil
}

func (a *App) extractArgs(args []string) []string {
	if len(args) == 1 {
		return []string{}
	}

	return args[1:]
}
