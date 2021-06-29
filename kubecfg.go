// command line tool to work with kubernetes configurations

package main

import (
	"os"

	"github.com/integrii/flaggy"
	"github.com/pherklotz/kubecfg/commands"
)

const VERSION = "1.1"

func main() {
	flaggy.SetName("kubecfg")
	flaggy.SetDescription("Small helper to manage kubernetes configurations and there contexts")
	flaggy.SetVersion(VERSION)
	flaggy.DebugMode = false

	cmdList := []commands.Command{
		commands.NewDeleteCommand(),
		commands.NewExportCommand(),
		commands.NewImportCommand(),
		commands.NewListCommand(),
		commands.NewRenameCommand(),
		commands.NewSwitchCommand(),
	}

	for _, cmd := range cmdList {
		description := cmd.GetCommand()
		flaggy.AttachSubcommand(description, 1)
	}
	flaggy.Parse()

	for _, cmd := range cmdList {
		if cmd.GetCommand().Used {
			cmd.Execute()
			os.Exit(0)
		}
	}
	flaggy.ShowHelpAndExit("")
}
