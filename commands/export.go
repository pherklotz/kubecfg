package commands

import (
	"fmt"
	"log"
	"regexp"
	"strconv"

	"github.com/pherklotz/kubecfg/common"

	"github.com/integrii/flaggy"
	k8s "k8s.io/client-go/tools/clientcmd/api/v1"
)

//NewExportCommand creates a new ExportCommand
func NewExportCommand() *ExportCommand {
	defaultPath, err := common.GetDefaultKubeconfigPath()
	if err != nil {
		log.Fatalln("Can not determine default path: ", err)
	}
	lc := ExportCommand{
		sourceFile: defaultPath,
	}

	cmd := flaggy.NewSubcommand("export")
	cmd.ShortName = "e"
	cmd.Description = "Exports a context into a new config file."
	cmd.AddPositionalValue(&lc.contextName, "context name", 1, false, "The name of the context to export.")
	cmd.String(&lc.sourceFile, "s", "source", "The optional path to the source kubeconfig file.")

	lc.command = cmd
	return &lc
}

//ExportCommand the list command struct
type ExportCommand struct {
	command     *flaggy.Subcommand
	sourceFile  string
	contextName string
}

// GetCommand returns the flaggy Subcommand to parse the command line
func (cmdArgs *ExportCommand) GetCommand() *flaggy.Subcommand {
	return cmdArgs.command
}

//Execute the export command
func (cmdArgs *ExportCommand) Execute(targetFile string) error {
	path := cmdArgs.sourceFile
	sourceConfig, err := common.ReadKubeConfigYaml(path)
	if err != nil {
		return fmt.Errorf("failed to load config from path '%s'.\nError: %v", path, err)
	}
	contextName := &cmdArgs.contextName

	if common.FileExists(targetFile) {
		reg, err := regexp.Compile("[^A-Za-z0-9]+")
		if err != nil {
			return err
		}
		cleanContextName := reg.ReplaceAllString(cmdArgs.contextName, "")
		targetFile = "kubeconf-" + cleanContextName
	}
	for i := 1; common.FileExists(targetFile); i++ {
		targetFile = targetFile + strconv.Itoa(i)
	}

	context, err := common.GetContextByName(sourceConfig, contextName)
	if err != nil {
		return fmt.Errorf("context with name '%s' not found", *contextName)
	}
	cluster, err := common.GetClusterByName(sourceConfig, &context.Context.Cluster)
	if err != nil {
		fmt.Printf("WARN: No associated cluster for context '%s' with name '%s' found.\n", *contextName, context.Context.Cluster)
	}
	user, err := common.GetUserByName(sourceConfig, &context.Context.AuthInfo)
	if err != nil {
		fmt.Printf("WARN: No associated user for context '%s' with name '%s' found.\n", *contextName, context.Context.AuthInfo)
	}

	targetConfig := k8s.Config{
		APIVersion:     sourceConfig.APIVersion,
		Kind:           sourceConfig.Kind,
		Contexts:       []k8s.NamedContext{*context},
		Clusters:       []k8s.NamedCluster{*cluster},
		AuthInfos:      []k8s.NamedAuthInfo{*user},
		CurrentContext: context.Name,
	}

	common.WriteKubeConfigYaml(targetFile, &targetConfig)
	fmt.Printf("Context '%s' exported to file '%s'", *contextName, targetFile)
	return nil
}
