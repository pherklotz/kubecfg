package commands

import (
	"bufio"
	"fmt"
	"kubecfg/common"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/integrii/flaggy"
	k8s "k8s.io/client-go/tools/clientcmd/api/v1"
)

//NewSwitchCommand a new SwitchCommand and return it.
func NewSwitchCommand() *SwitchCommand {
	sc := SwitchCommand{}

	cmd := flaggy.NewSubcommand("switch")
	cmd.ShortName = "s"
	cmd.Description = "Switches the active context. Starts interactive mode if no context is given."
	cmd.AddPositionalValue(&sc.contextName, "context", 1, false, "The name of the context to activate.")
	sc.command = cmd
	return &sc
}

//SwitchCommand the activate command struct
type SwitchCommand struct {
	command     *flaggy.Subcommand
	contextName string
}

// GetCommand returns the flaggy Subcommand to parse the command line
func (cmdArgs *SwitchCommand) GetCommand() *flaggy.Subcommand {
	return cmdArgs.command
}

//Execute the list command
func (cmdArgs *SwitchCommand) Execute() {
	path, err := common.GetDefaultKubeconfigPath()
	if err != nil {
		log.Fatalf("Failed to load default config path.\nError: %v\n", err)
	}
	config, err := common.ReadKubeConfigYaml(path)
	if err != nil {
		log.Fatalf("Failed to load config from path '%s'.\nError: %v\n", path, err)
	}

	var context k8s.NamedContext
	if cmdArgs.contextName != "" {
		context = getContextByName(config.Contexts, &cmdArgs.contextName)
	} else {
		context = getTargetConfigWithInteractiveMode(config)

	}

	config.CurrentContext = context.Name

	common.WriteKubeConfigYaml(path, config)
	fmt.Printf("Activate context: %s\n", config.CurrentContext)
}

func getContextByName(contexts []k8s.NamedContext, contextName *string) k8s.NamedContext {
	for _, context := range contexts {
		if context.Name == *contextName {
			return context
		}
	}
	log.Fatalf("Unknown context name '%s'. Use 'kubecfg list' to see available contexts.\n", *contextName)
	return k8s.NamedContext{}
}

func getTargetConfigWithInteractiveMode(config *k8s.Config) k8s.NamedContext {
	contexts := config.Contexts
	for index, context := range contexts {
		activeIndicator := " "
		if context.Name == config.CurrentContext {
			activeIndicator = "*"
		}
		fmt.Printf("[%d] %s %s\n", index, activeIndicator, context.Name)
	}
	fmt.Printf("Enter number of context to activate: ")
	reader := bufio.NewReader(os.Stdin)
	optionString, err := reader.ReadString('\n')
	if err != nil {
		log.Fatal("Can not read selected option.", err)
	}
	selectedOption, err := strconv.Atoi(strings.TrimSpace(optionString))
	if err != nil {
		log.Fatal("Please enter a valid option number", err)
	}
	if selectedOption < 0 || selectedOption >= len(contexts) {
		log.Fatalln("Option is not valid. Please enter a valid option number.")
	}
	return contexts[selectedOption]
}
