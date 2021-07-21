package commands

import (
	"fmt"
	"log"
	"sort"

	"github.com/pherklotz/kubecfg/common"

	"github.com/integrii/flaggy"
)

//NewListCommand creates a new ListCommand
func NewListCommand() *ListCommand {
	lc := ListCommand{}

	cmd := flaggy.NewSubcommand("list")
	cmd.ShortName = "l"
	cmd.Description = "Lists all contexts in the config file."

	lc.command = cmd
	return &lc
}

//ListCommand the list command struct
type ListCommand struct {
	command *flaggy.Subcommand
}

// GetCommand returns the flaggy Subcommand to parse the command line
func (cmdArgs *ListCommand) GetCommand() *flaggy.Subcommand {
	return cmdArgs.command
}

//Execute the list command
func (cmdArgs *ListCommand) Execute(path string) error {
	config, err := common.ReadKubeConfigYaml(path)
	if err != nil {
		log.Fatalf("Failed to load config from path '%s'.\nError: %v\n", path, err)
	}
	contexts := config.Contexts

	sort.Slice(contexts, func(i, j int) bool {
		return contexts[i].Name < contexts[j].Name
	})

	for _, context := range contexts {
		activeIndicator := " "
		if context.Name == config.CurrentContext {
			activeIndicator = "*"
		}
		fmt.Printf("%s %s\n", activeIndicator, context.Name)
	}
	return nil
}
