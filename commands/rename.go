package commands

import (
	"fmt"

	"github.com/pherklotz/kubecfg/common"

	"github.com/integrii/flaggy"
	k8s "k8s.io/client-go/tools/clientcmd/api/v1"
)

//NewRenameCommand a new RenameCommand and return it.
func NewRenameCommand() *RenameCommand {
	rc := RenameCommand{}

	cmd := flaggy.NewSubcommand("rename")
	cmd.ShortName = "rn"
	cmd.Description = "Renames a context."
	cmd.AddPositionalValue(&rc.from, "from", 1, true, "The old context name.")
	cmd.AddPositionalValue(&rc.to, "to", 2, true, "The new context name.")
	rc.command = cmd
	return &rc
}

//RenameCommand the rename command struct
type RenameCommand struct {
	command *flaggy.Subcommand
	from    string
	to      string
}

// GetCommand returns the flaggy Subcommand to parse the command line
func (cmdArgs *RenameCommand) GetCommand() *flaggy.Subcommand {
	return cmdArgs.command
}

//Execute the rename command
func (cmdArgs *RenameCommand) Execute(path string) error {
	config, err := common.ReadKubeConfigYaml(path)
	if err != nil {
		return err
	}
	newConfig, err := rename(config, &cmdArgs.from, &cmdArgs.to)
	if err != nil {
		return err
	}
	common.WriteKubeConfigYaml(path, newConfig)
	return nil
}

// Rename renames the context from oldName to newName
func rename(config *k8s.Config, oldName *string, newName *string) (*k8s.Config, error) {
	if *newName == "" {
		return nil, fmt.Errorf("Please use '-to <name>' to specify the new name.")
	}

	context, err := common.GetContextByName(config, oldName)
	if err != nil {
		return nil, err
	}

	_, err = common.GetContextByName(config, newName)
	if err == nil {
		return nil, fmt.Errorf("There is already a context with the name: '%s'\n", *newName)
	}

	cluster, err := common.GetClusterByName(config, &context.Context.Cluster)
	if err != nil {
		return nil, fmt.Errorf("No cluster found for context '%s'. %v\n", context.Context.Cluster, err)
	}
	user, err := common.GetUserByName(config, &context.Context.AuthInfo)
	if err != nil {
		return nil, fmt.Errorf("No user found for context '%s'. %v\n", context.Context.AuthInfo, err)
	}
	cluster.Name = *newName
	user.Name = *newName
	context.Name = *newName
	context.Context.Cluster = *newName
	context.Context.AuthInfo = *newName
	if config.CurrentContext == *oldName {
		config.CurrentContext = *newName
	}
	return config, nil
}
