package commands

import (
	"fmt"
	"log"
	"regexp"

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
		targetFile: "",
	}

	cmd := flaggy.NewSubcommand("export")
	cmd.ShortName = "e"
	cmd.Description = "Exports a context into a new config file."
	cmd.AddPositionalValue(&lc.contextName, "context name", 1, false, "The name of the context to export.")
	cmd.String(&lc.sourceFile, "s", "source", "The optional path to the source kubeconfig file.")
	cmd.String(&lc.targetFile, "t", "target", "The optional path to the new kubeconfig file.")

	lc.command = cmd
	return &lc
}

//ExportCommand the list command struct
type ExportCommand struct {
	command     *flaggy.Subcommand
	sourceFile  string
	targetFile  string
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
		log.Fatalf("Failed to load config from path '%s'.\nError: %v\n", path, err)
	}
	contextName := &cmdArgs.contextName

	defaultKubecfgFile, err := common.GetDefaultKubeconfigPath()
	if err != nil {
		log.Fatalf("Failed to load default kubeconfig path.\nError: %v\n", err)
	}

	// check if is the default value, if so create a new file name
	if targetFile == defaultKubecfgFile {
		reg, err := regexp.Compile("[^A-Za-z0-9]+")
		if err != nil {
			log.Fatal(err)
		}
		cleanContextName := reg.ReplaceAllString(cmdArgs.contextName, "")
		targetFile = "kubeconf-" + cleanContextName
	}

	if common.FileExists(targetFile) {
		log.Fatalf("Target file '%s' exists already.\n", targetFile)
	}

	context, err := common.GetContextByName(sourceConfig, contextName)
	if err != nil {
		log.Fatalf("Context with name '%s' not found.\n", *contextName)
	}
	cluster, err := common.GetClusterByName(sourceConfig, &context.Context.Cluster)
	if err != nil {
		log.Printf("WARN: No associated cluster for context '%s' with name '%s' found.\n", *contextName, context.Context.Cluster)
	}
	user, err := common.GetUserByName(sourceConfig, &context.Context.AuthInfo)
	if err != nil {
		log.Printf("WARN: No associated user for context '%s' with name '%s' found.\n", *contextName, context.Context.AuthInfo)
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
