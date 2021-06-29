package commands

import (
	"log"

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
func (cmdArgs *RenameCommand) Execute() {
	path, err := common.GetDefaultKubeconfigPath()
	if err != nil {
		log.Fatalf("Failed to load default config path.\nError: %v\n", err)
	}
	config, err := common.ReadKubeConfigYaml(path)
	if err != nil {
		log.Fatalf("Failed to load config from path '%s'.\nError: %v\n", path, err)
	}
	newConfig := Rename(config, &cmdArgs.from, &cmdArgs.to)

	common.WriteKubeConfigYaml(path, newConfig)
}

// Rename renames the context from oldName to newName
func Rename(config *k8s.Config, oldName *string, newName *string) *k8s.Config {
	if *newName == "" {
		log.Fatalln("Please use '-to <name>' to specify the new name.")
	}

	context, err := common.GetContextByName(config, oldName)
	if err != nil {
		log.Fatalln(err)
	}

	_, err = common.GetContextByName(config, newName)
	if err == nil {
		log.Fatalf("There is already a context with the name: '%s'\n", *newName)
	}

	cluster, err := common.GetClusterByName(config, &context.Context.Cluster)
	if err != nil {
		log.Fatalf("No cluster found for context '%s'. %v\n", context.Context.Cluster, err)
	}
	user, err := common.GetUserByName(config, &context.Context.AuthInfo)
	if err != nil {
		log.Fatalf("No user found for context '%s'. %v\n", context.Context.AuthInfo, err)
	}
	cluster.Name = *newName
	user.Name = *newName
	context.Name = *newName
	context.Context.Cluster = *newName
	context.Context.AuthInfo = *newName
	if config.CurrentContext == *oldName {
		config.CurrentContext = *newName
	}
	return config
}
