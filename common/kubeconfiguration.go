package common

import (
	"fmt"

	k8s "k8s.io/client-go/tools/clientcmd/api/v1"
)

// GetContextByName returns the first context with the given name
func GetContextByName(config *k8s.Config, contextName *string) (*k8s.NamedContext, error) {
	for i := range config.Contexts {
		context := &config.Contexts[i]
		if context.Name == *contextName {
			return context, nil
		}
	}
	return nil, fmt.Errorf("Can not find context with name '%s'", *contextName)
}

//GetClusterByName returns the cluster with the given name or an error
func GetClusterByName(config *k8s.Config, clusterName *string) (*k8s.NamedCluster, error) {
	for i := range config.Clusters {
		cluster := &config.Clusters[i]
		if cluster.Name == *clusterName {
			return cluster, nil
		}
	}
	return nil, fmt.Errorf("Can not find cluster with name '%s'", *clusterName)
}

//GetUserByName returns the user with the given name or an error
func GetUserByName(config *k8s.Config, userName *string) (*k8s.NamedAuthInfo, error) {
	for i := range config.AuthInfos {
		user := &config.AuthInfos[i]
		if user.Name == *userName {
			return user, nil
		}
	}
	return nil, fmt.Errorf("Can not find user with name '%s'", *userName)
}
