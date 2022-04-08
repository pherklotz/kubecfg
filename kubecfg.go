// command line tool to work with kubernetes configurations
package main

import (
	"log"
	"os"

	"github.com/integrii/flaggy"
	"github.com/pherklotz/kubecfg/commands"
	"github.com/pherklotz/kubecfg/common"
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
		commands.NewViewCommand(),
	}

	for _, cmd := range cmdList {
		commandDescription := cmd.GetCommand()
		flaggy.AttachSubcommand(commandDescription, 1)
	}

	targetFile, err := common.GetDefaultKubeconfigPath()
	if err != nil {
		log.Fatalf("Failed to load default kubeconfig path.\nError: %v\n", err)
	}
	flaggy.String(&targetFile, "t", "target", "The target file for the specified command. The file must be a valid kubeconfig file.")
	flaggy.Parse()

	for _, cmd := range cmdList {
		if cmd.GetCommand().Used {
			err = cmd.Execute(targetFile)
			if err != nil {
				log.Fatal(err)
			}
			os.Exit(0)
		}
	}
	flaggy.ShowHelpAndExit("")
}
