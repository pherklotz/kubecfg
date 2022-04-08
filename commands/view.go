package commands

import (
	"fmt"

	"github.com/integrii/flaggy"
	"github.com/pherklotz/kubecfg/common"
	k8s "k8s.io/client-go/tools/clientcmd/api/v1"
)

//NewViewCommand creates a new ViewCommand
func NewViewCommand() *ViewCommand {
	vc := ViewCommand{}

	cmd := flaggy.NewSubcommand("view")
	cmd.ShortName = "v"
	cmd.Description = "Prints the details of a context."
	cmd.AddPositionalValue(&vc.context, "context", 1, false, "The name of the context to print. Default: activated context")

	vc.command = cmd
	return &vc
}

//ViewCommand the view command struct
type ViewCommand struct {
	command *flaggy.Subcommand
	context string
}

// GetCommand returns the flaggy Subcommand to parse the command line
func (cmdArgs *ViewCommand) GetCommand() *flaggy.Subcommand {
	return cmdArgs.command
}

//Execute the view command
func (cmdArgs *ViewCommand) Execute(path string) error {
	config, err := common.ReadKubeConfigYaml(path)
	if cmdArgs.context == "" {
		cmdArgs.context = config.CurrentContext
	}

	if err != nil {
		return fmt.Errorf("Failed to load config from path '%s'.\nError: %v\n", path, err)
	}
	context, err := common.GetContextByName(config, &cmdArgs.context)
	if err != nil {
		return err
	}
	cluster, err := common.GetClusterByName(config, &context.Context.Cluster)
	if err != nil {
		return err
	}
	user, err := common.GetUserByName(config, &context.Context.Cluster)
	if err != nil {
		return err
	}

	newConfig := k8s.Config{
		Contexts:       []k8s.NamedContext{*context},
		Clusters:       []k8s.NamedCluster{*cluster},
		AuthInfos:      []k8s.NamedAuthInfo{*user},
		CurrentContext: context.Name,
	}
	configStr, err := common.ConfigToString(&newConfig)
	if err != nil {
		return err
	}
	fmt.Print(configStr)
	return nil
}
