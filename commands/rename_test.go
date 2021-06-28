package commands

import (
	"testing"

	"github.com/stretchr/testify/assert"
	k8s "k8s.io/client-go/tools/clientcmd/api/v1"
)

func TestRenameCommand_Rename(t *testing.T) {
	oldName := "old_name"
	context := k8s.NamedContext{Name: oldName}
	context.Context.Cluster = oldName
	context.Context.AuthInfo = oldName
	cluster := k8s.NamedCluster{Name: oldName}
	user := k8s.NamedAuthInfo{Name: oldName}

	kubeCfg := k8s.Config{}
	kubeCfg.Contexts = make([]k8s.NamedContext, 1)
	kubeCfg.Contexts[0] = context
	kubeCfg.Clusters = make([]k8s.NamedCluster, 1)
	kubeCfg.Clusters[0] = cluster
	kubeCfg.AuthInfos = make([]k8s.NamedAuthInfo, 1)
	kubeCfg.AuthInfos[0] = user
	kubeCfg.CurrentContext = oldName

	newName := "new_name"
	changedKubeconfig := Rename(&kubeCfg, &oldName, &newName)

	assert.Equal(t, newName, changedKubeconfig.Clusters[0].Name, "KubeConf.Clusters[0].Name")
	assert.Equal(t, newName, changedKubeconfig.Contexts[0].Name, "KubeConf.Contexts[0].Name")
	assert.Equal(t, newName, changedKubeconfig.Contexts[0].Context.Cluster, "KubeConf.Contexts[0].Context.Cluster")
	assert.Equal(t, newName, changedKubeconfig.Contexts[0].Context.User, "KubeConf.Contexts[0].Context.User")
	assert.Equal(t, newName, changedKubeconfig.Users[0].Name, "KubeConf.Users[0].Name")
	assert.Equal(t, newName, changedKubeconfig.CurrentContext, "KubeConf.CurrentContext")
}
