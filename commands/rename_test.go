package commands

import (
	"fmt"
	"os"
	"testing"

	"github.com/pherklotz/kubecfg/common"
	"github.com/stretchr/testify/assert"
	k8s "k8s.io/client-go/tools/clientcmd/api/v1"
)

func createKubeconf() k8s.Config {
	oldName := "old_name"
	context := k8s.NamedContext{Name: oldName}
	context.Context.Cluster = oldName
	context.Context.AuthInfo = oldName
	cluster := k8s.NamedCluster{Name: oldName}
	user := k8s.NamedAuthInfo{Name: oldName}

	kubeCfg := k8s.Config{}
	kubeCfg.Contexts = make([]k8s.NamedContext, 2)
	kubeCfg.Contexts[0] = context
	kubeCfg.Contexts[1] = k8s.NamedContext{Name: "foo"}
	kubeCfg.Clusters = make([]k8s.NamedCluster, 1)
	kubeCfg.Clusters[0] = cluster
	kubeCfg.AuthInfos = make([]k8s.NamedAuthInfo, 1)
	kubeCfg.AuthInfos[0] = user
	kubeCfg.CurrentContext = oldName
	return kubeCfg
}

func TestRenameCommand(t *testing.T) {
	tempDir := t.TempDir()

	kubeCfg := createKubeconf()

	var tests = []struct {
		cmd          RenameCommand
		expectedName string
		expectError  bool
	}{
		{
			RenameCommand{
				from: kubeCfg.Contexts[0].Name,
				to:   "new_name",
			},
			"new_name",
			false,
		},
		{
			RenameCommand{
				from: kubeCfg.Contexts[0].Name,
				to:   "",
			},
			kubeCfg.Contexts[0].Name,
			true,
		},
		{
			RenameCommand{
				from: "UNKNOWN",
				to:   "new_name",
			},
			kubeCfg.Contexts[0].Name,
			true,
		},
		{
			RenameCommand{
				from: kubeCfg.Contexts[0].Name,
				to:   kubeCfg.Contexts[1].Name,
			},
			kubeCfg.Contexts[0].Name,
			true,
		},
		{
			RenameCommand{
				from: kubeCfg.Contexts[1].Name,
				to:   "new_name",
			},
			kubeCfg.Contexts[0].Name,
			true,
		},
	}
	for index, tt := range tests {
		testname := fmt.Sprintf("Expect '%s' on %d call of rename command", tt.expectedName, index+1)
		t.Run(testname, func(t *testing.T) {
			kubeCfgTarget := tempDir + string(os.PathSeparator) + fmt.Sprint("kubecfg_", index)
			common.WriteKubeConfigYaml(kubeCfgTarget, &kubeCfg)

			err := tt.cmd.Execute(kubeCfgTarget)
			if tt.expectError {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
			}

			changedKubeconfig, err := common.ReadKubeConfigYaml(kubeCfgTarget)
			assert.Nil(t, err)

			assert.Equal(t, tt.expectedName, changedKubeconfig.Clusters[0].Name, "KubeConf.Clusters[0].Name")
			assert.Equal(t, tt.expectedName, changedKubeconfig.Contexts[0].Name, "KubeConf.Contexts[0].Name")
			assert.Equal(t, tt.expectedName, changedKubeconfig.Contexts[0].Context.Cluster, "KubeConf.Contexts[0].Context.Cluster")
			assert.Equal(t, tt.expectedName, changedKubeconfig.Contexts[0].Context.AuthInfo, "KubeConf.Contexts[0].Context.AuthInfos")
			assert.Equal(t, tt.expectedName, changedKubeconfig.AuthInfos[0].Name, "KubeConf.AuthInfos[0].Name")
			assert.Equal(t, tt.expectedName, changedKubeconfig.CurrentContext, "KubeConf.CurrentContext")
		})
	}
}
