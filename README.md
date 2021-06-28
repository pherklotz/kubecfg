# kubecfg v1.0
Command line tool for easy management of kubernetes configurations. This tool might help you if you often have to change between several kubernetes configurations. Therefor it will mostly modify the default kubernetes config file (~/.kube/config). So you switch the kubernetes context globally and not only for this shell!

Currently the following features are integrated:

- list all contexts in the default kubeconfig 
- switch active context globally
- import a context from one file into the default config file
- export a context from the default config file into a separate file
- delete a context
## Usage

kubecfg contains a help that should cover the most commands and flags. The help works also for specific commands.

Examples:
```sh
# Show the general help
kubecfg help
# Show the help for a specific command
kubecfg help <command>
# Lists all known contexts
kubecfg list
# Switch to a context in interactive mode
kubecfg switch
# Switch direct to a context
kubecfg switch <context name>
```
# Installation

