package commands

import (
	"fmt"
	"log"
	"path/filepath"

	"github.com/pherklotz/kubecfg/common"

	"github.com/integrii/flaggy"
	k8s "k8s.io/client-go/tools/clientcmd/api/v1"
)

//NewImportCommand creates a new Importommand
func NewImportCommand() *ImportCommand {
	ac := ImportCommand{activate: false, name: "", filePattern: "kubeconf-"}

	cmd := flaggy.NewSubcommand("import")
	cmd.ShortName = "i"
	cmd.Description = "Adds one or more contexts of found kubeconf files to the home config file."
	cmd.AddPositionalValue(&ac.filePattern, "filePattern", 1, false, "The glob pattern to search for kubeconfig files.")
	cmd.Bool(&ac.activate, "s", "switch", "Switch to the last added context.")
	cmd.String(&ac.name, "n", "name", "Name of the new context. If the name is not unique in the target kubeconfig file a suffix will be attached.")

	ac.command = cmd
	return &ac
}

// ImportCommand to add a kubeconf
type ImportCommand struct {
	command     *flaggy.Subcommand
	filePattern string
	activate    bool
	name        string
}

// GetCommand returns the flaggy Subcommand to parse the command line
func (cmdArgs *ImportCommand) GetCommand() *flaggy.Subcommand {
	return cmdArgs.command
}

//Execute the import command
func (cmdArgs *ImportCommand) Execute() {
	var ctxNameProvider common.ContextNameProvider
	if cmdArgs.name == "" {
		ctxNameProvider = &common.RandomNameProvider{}
	} else {
		ctxNameProvider = &common.SeedNameProvider{Seed: cmdArgs.name}
	}

	addConfig(&cmdArgs.filePattern, &cmdArgs.activate, ctxNameProvider)
}

// addConfig searchs for kubeconfig files and add them to the default config
func addConfig(sourcePattern *string, activate *bool, ctxNameProvider common.ContextNameProvider) {
	configFiles, err := filepath.Glob(*sourcePattern)
	if err != nil {
		log.Fatalln("Wrong source file pattern '", *sourcePattern, "':", err)
	}
	foundFilesCount := len(configFiles)
	if foundFilesCount == 0 {
		log.Println("Found no matching files to source pattern: ", *sourcePattern)
	}

	targetConfFile, err := common.GetDefaultKubeconfigPath()
	if err != nil {
		log.Fatalln("can not read user home: ", err)
	}
	targetConf, err := common.ReadKubeConfigYaml(targetConfFile)
	if err != nil {
		log.Fatalln("Can not parse kubeconfig yaml in file '", targetConfFile, "':", err)
	}

	for i := 0; i < foundFilesCount; i++ {
		sourceFile := configFiles[i]
		sourceConf, err := common.ReadKubeConfigYaml(sourceFile)
		if err != nil {
			log.Fatalln("Can not parse kubeconfig yaml in file '", sourceFile, "':", err)
		}
		fmt.Printf("Add kubeconf from '%s' to '%s'\n", sourceFile, targetConfFile)

		for j := 0; j < len(sourceConf.Contexts); j++ {
			context := sourceConf.Contexts[j]
			oldName := context.Name
			cluster, err := findCluster(sourceConf.Clusters, context.Context.Cluster)
			if err != nil {
				log.Fatalf("No cluster with name '%s' in file '%s'. Error: %v", context.Context.Cluster, sourceFile, err)
			}

			user, err := findUser(sourceConf.AuthInfos, context.Context.AuthInfo)
			if err != nil {
				log.Fatalf("No user with name '%s' in file '%s'. Error: %v", context.Context.AuthInfo, sourceFile, err)
			}

			newName := getUniqueName(sourceConf.Contexts, ctxNameProvider)
			context.Name = newName
			context.Context.Cluster = newName
			context.Context.AuthInfo = newName

			cluster.Name = newName
			user.Name = newName

			targetConf.Contexts = append(targetConf.Contexts, context)
			targetConf.Clusters = append(targetConf.Clusters, cluster)
			targetConf.AuthInfos = append(targetConf.AuthInfos, user)
			fmt.Printf("\tAdded context '%s' with new name '%s'\n", oldName, newName)
			if *activate {
				targetConf.CurrentContext = newName
			}
		}
	}
	if *activate {
		fmt.Printf("Activated context: %s\n", targetConf.CurrentContext)
	}
	err = common.CopyFile(targetConfFile, targetConfFile+".bak")
	if err != nil {
		log.Fatalf("Could no create backup of target '%s': %v", targetConfFile, err)
	}
	common.WriteKubeConfigYaml(targetConfFile, targetConf)
}

func getUniqueName(contexts []k8s.NamedContext, nameProvider common.ContextNameProvider) string {
	newName := nameProvider.GetName()
	for {
		if !containsName(contexts, newName) {
			return newName
		}
	}
}

func containsName(contexts []k8s.NamedContext, name string) bool {
	for _, ctx := range contexts {
		if ctx.Name == name || ctx.Context.Cluster == name || ctx.Context.AuthInfo == name {
			return true
		}
	}
	return false
}

func findUser(users []k8s.NamedAuthInfo, name string) (k8s.NamedAuthInfo, error) {
	for i := 0; i < len(users); i++ {
		if users[i].Name == name {
			return users[i], nil
		}
	}
	return k8s.NamedAuthInfo{}, fmt.Errorf("Did not find an user with name '%s'", name)
}

func findCluster(clusters []k8s.NamedCluster, name string) (k8s.NamedCluster, error) {
	for i := 0; i < len(clusters); i++ {
		if clusters[i].Name == name {
			return clusters[i], nil
		}
	}
	return k8s.NamedCluster{}, fmt.Errorf("Did not find a cluster with name '%s'", name)
}
