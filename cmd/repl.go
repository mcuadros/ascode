package cmd

import "github.com/jessevdk/go-flags"

// Command descriptions used in the flags.Parser.AddCommand.
const (
	REPLCmdShortDescription = "Run as interactive shell."
	REPLCmdLongDescription  = REPLCmdShortDescription + "\n\n" +
		"The REPL shell provides the same capabilities as the regular `run`\n" +
		"command."
)

// REPLCmd implements the command `repl`.
type REPLCmd struct {
	commonCmd
}

// Execute honors the flags.Commander interface.
func (c *REPLCmd) Execute(args []string) error {
	c.init()
	c.runtime.REPL()

	return nil
}

var _ flags.Commander = &REPLCmd{}
