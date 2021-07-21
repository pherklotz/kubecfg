package commands

import "github.com/integrii/flaggy"

// Command defines a command line command.
type Command interface {
	GetCommand() *flaggy.Subcommand

	Execute(targetFilePath string) error
}
