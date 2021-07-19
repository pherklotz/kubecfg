package commands

import (
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
func (cmdArgs *DeleteCommand) Execute(path string) {
	config, err := common.ReadKubeConfigYaml(path)
	if err != nil {
		log.Fatalf("Failed to load config from path: %s\nError: %v\n", path, err)
	}
	contextName := cmdArgs.context
	if contextName == config.CurrentContext {
		log.Fatalf("Can not delete the current context '%s'. Change the context first.\n", contextName)
	}

	contextToDelete, err := common.GetContextByName(config, &contextName)
	if err != nil {
		log.Fatalf("Context with name '%s' not found\n", contextName)
	}

	config.Contexts = deleteContext(config.Contexts, contextToDelete.Name)
	config.AuthInfos = deleteUser(config.AuthInfos, contextToDelete.Context.AuthInfo)
	config.Clusters = deleteCluster(config.Clusters, contextToDelete.Context.Cluster)

	common.WriteKubeConfigYaml(path, config)

	log.Printf("Deleted context with name '%s' and associated clusters and users.\n", contextToDelete.Name)
}

func deleteContext(contexts []k8s.NamedContext, contextNameToDelete string) []k8s.NamedContext {
	for i, ctx := range contexts {
		if ctx.Name == contextNameToDelete {
			if i+1 < len(contexts) {
				contexts = append(contexts[:i], contexts[i+1:]...)
			} else {
				contexts = append(contexts[:i])
			}
			break
		}
	}
	return contexts
}

func deleteUser(authInfos []k8s.NamedAuthInfo, toDelete string) []k8s.NamedAuthInfo {
	for i, authInfo := range authInfos {
		if authInfo.Name == toDelete {
			if i+1 < len(authInfos) {
				authInfos = append(authInfos[:i], authInfos[i+1:]...)
			} else {
				authInfos = append(authInfos[:i])
			}
			break
		}
	}
	return authInfos
}

func deleteCluster(clusters []k8s.NamedCluster, toDelete string) []k8s.NamedCluster {
	for i, cluster := range clusters {
		if cluster.Name == toDelete {
			if i+1 < len(clusters) {
				clusters = append(clusters[:i], clusters[i+1:]...)
			} else {
				clusters = append(clusters[:i])
			}
			break
		}
	}
	return clusters
}
