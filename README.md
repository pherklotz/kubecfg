# kubecfg
Command line tool for easy management of kubernetes configurations. This tool might help you if you often have to change between several kubernetes configurations. Therefor it will mostly modify the default kubernetes config file (~/.kube/config). So you switch the kubernetes context globally and not only for this shell!

Currently the following features are integrated:

- list all contexts in the default kubeconfig 
- switch active context globally
- import a context from one file into the default config file
- export a context from the default config file into a separate file
- rename a context
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
# Switch directly to a context
kubecfg switch <context name>
# Rename a context
kubecfg rename <old context name> <new context name> 
# Import a context file into the default config. Glob pattern is allowed.
kubecfg import <path to kubeconfig file> 
# Exports a context from the default config into a new file
kubecfg export <context name> 
# Deletes a context and all associated data from the default config
kubecfg delete <context name> 
```
# Installation
- Download the file for your system from the release section
- Rename the file to kubecfg
- Add the directory to your PATH variable 

# Disclaimer
Kubecfg is open source and you use it at your own risk.
We are not responsible for any data loses. Be careful and backup your data regularly.

# License
The project is under the Apache 2.0 License
