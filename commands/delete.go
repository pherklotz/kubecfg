package commands

import (
	"fmt"
	"log"

	"github.com/integrii/flaggy"
	"github.com/pherklotz/kubecfg/common"
	k8s "k8s.io/client-go/tools/clientcmd/api/v1"
)

//NewDeleteCommand creates a new DeleteCommand
func NewDeleteCommand() *DeleteCommand {
	lc := DeleteCommand{}

	cmd := flaggy.NewSubcommand("delete")
	cmd.ShortName = "del"
	cmd.Description = "Deletes the context from the config file. Fails if it is the active context."
	cmd.AddPositionalValue(&lc.context, "context", 1, true, "The name of the context to delete.")

	lc.command = cmd
	return &lc
}

//DeleteCommand the command struct
type DeleteCommand struct {
	command *flaggy.Subcommand
	context string
}

// GetCommand returns the flaggy Subcommand to parse the command line
func (cmdArgs *DeleteCommand) GetCommand() *flaggy.Subcommand {
	return cmdArgs.command
}

//Execute the delete command
func (cmdArgs *DeleteCommand) Execute(path string) error {
	config, err := common.ReadKubeConfigYaml(path)
	if err != nil {
		return fmt.Errorf("failed to load config from path: %s\nError: %v", path, err)
	}
	contextName := cmdArgs.context
	if contextName == config.CurrentContext {
		return fmt.Errorf("can not delete the current context '%s'. Change the context first", contextName)
	}

	contextToDelete, err := common.GetContextByName(config, &contextName)
	if err != nil {
		return err
	}

	config.Contexts = filterContexts(config.Contexts, contextToDelete.Name)
	config.AuthInfos = filterAuthInfos(config.AuthInfos, contextToDelete.Context.AuthInfo)
	config.Clusters = filterClusters(config.Clusters, contextToDelete.Context.Cluster)

	common.WriteKubeConfigYaml(path, config)

	log.Printf("Deleted context with name '%s' and associated clusters and users.", contextToDelete.Name)
	return nil
}

func filterAuthInfos(list []k8s.NamedAuthInfo, toDelete string) []k8s.NamedAuthInfo {
	filteredList := make([]k8s.NamedAuthInfo, 0, len(list))
	for _, item := range list {
		if item.Name != toDelete {
			filteredList = append(filteredList, item)
		}
	}
	return filteredList
}

func filterClusters(list []k8s.NamedCluster, toDelete string) []k8s.NamedCluster {
	filteredList := make([]k8s.NamedCluster, 0, len(list))
	for _, item := range list {
		if item.Name != toDelete {
			filteredList = append(filteredList, item)
		}
	}
	return filteredList
}

func filterContexts(list []k8s.NamedContext, toDelete string) []k8s.NamedContext {
	filteredList := make([]k8s.NamedContext, 0, len(list))
	for _, item := range list {
		if item.Name != toDelete {
			filteredList = append(filteredList, item)
		}
	}
	return filteredList
}
