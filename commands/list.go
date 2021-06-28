package commands

import (
	"fmt"
	"kubecfg/common"
	"log"
	"sort"

	"github.com/integrii/flaggy"
)

//NewListCommand creates a new ListCommand
func NewListCommand() *ListCommand {
	defaultPath, err := common.GetDefaultKubeconfigPath()
	if err != nil {
		log.Fatalln("Can not determine default path: ", err)
	}
	lc := ListCommand{sourceFile: defaultPath}

	cmd := flaggy.NewSubcommand("list")
	cmd.ShortName = "l"
	cmd.Description = "Lists all contexts in the config file."
	cmd.AddPositionalValue(&lc.sourceFile, "source", 1, false, "The optional path to the kubeconfig file.")

	lc.command = cmd
	return &lc
}

//ListCommand the list command struct
type ListCommand struct {
	command    *flaggy.Subcommand
	sourceFile string
}

// GetCommand returns the flaggy Subcommand to parse the command line
func (cmdArgs *ListCommand) GetCommand() *flaggy.Subcommand {
	return cmdArgs.command
}

//Execute the list command
func (cmdArgs *ListCommand) Execute() {
	path := cmdArgs.sourceFile
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
}
