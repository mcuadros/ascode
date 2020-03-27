package cmd

const (
	REPLCmdShortDescription = "Run as interactive shell."
	REPLCmdLongDescription  = REPLCmdShortDescription + "\n\n" +
		"The REPL shell provides the same capabilities as the regular `run`\n" +
		"command."
)

type REPLCmd struct {
	commonCmd
}

func (c *REPLCmd) Execute(args []string) error {
	c.init()
	c.runtime.REPL()

	return nil
}
