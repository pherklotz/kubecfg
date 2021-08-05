package commands

import (
	"fmt"
	"os"
	"testing"

	"github.com/pherklotz/kubecfg/common"
	"github.com/stretchr/testify/assert"
	k8s "k8s.io/client-go/tools/clientcmd/api/v1"
)

func createTestConfig(ctxName string) k8s.Config {
	context := k8s.NamedContext{Name: ctxName, Context: k8s.Context{Cluster: ctxName, AuthInfo: ctxName}}

	kubeCfg := k8s.Config{}
	kubeCfg.Contexts = make([]k8s.NamedContext, 2)
	kubeCfg.Contexts[0] = context
	kubeCfg.Contexts[1] = k8s.NamedContext{Name: "foo"}
	kubeCfg.Clusters = make([]k8s.NamedCluster, 1)
	kubeCfg.Clusters[0] = k8s.NamedCluster{Name: ctxName}
	kubeCfg.AuthInfos = make([]k8s.NamedAuthInfo, 1)
	kubeCfg.AuthInfos[0] = k8s.NamedAuthInfo{Name: ctxName}
	kubeCfg.CurrentContext = "activeCtx"
	return kubeCfg
}

func TestDeleteCommand(t *testing.T) {
	ctxName := "name"
	kubeCfg := createTestConfig(ctxName)

	tempDir := t.TempDir()

	var tests = []struct {
		cmd         DeleteCommand
		expectError bool
	}{
		{
			DeleteCommand{
				context: ctxName,
			},
			false,
		},
		{
			DeleteCommand{
				context: kubeCfg.CurrentContext,
			},
			true,
		},
		{
			DeleteCommand{
				context: "UNKNOWN",
			},
			true,
		},
	}
	for index, tt := range tests {
		expectErrorStr := "no"
		if tt.expectError {
			expectErrorStr = "an"
		}
		testname := fmt.Sprintf("Expect %s error on %d call of delete command. Try to delete context '%s'", expectErrorStr, index+1, tt.cmd.context)
		t.Run(testname, func(t *testing.T) {
			kubeCfgTarget := tempDir + string(os.PathSeparator) + fmt.Sprint("kubecfg_", index)
			common.WriteKubeConfigYaml(kubeCfgTarget, &kubeCfg)

			err := tt.cmd.Execute(kubeCfgTarget)

			changedKubeconfig, readErr := common.ReadKubeConfigYaml(kubeCfgTarget)
			assert.Nil(t, readErr)

			if tt.expectError {
				assert.NotNil(t, err)
				assert.EqualValues(t, kubeCfg, *changedKubeconfig)
			} else {
				assert.Nil(t, err)
				assert.Equal(t, 1, len(changedKubeconfig.Contexts), "len KubeConf.Contexts")
				assert.Equal(t, "foo", changedKubeconfig.Contexts[0].Name, "Context name")
				assert.Equal(t, 0, len(changedKubeconfig.Clusters), "len KubeConf.Clusters")
				assert.Equal(t, 0, len(changedKubeconfig.AuthInfos), "len KubeConf.AuthInfos")
				assert.Equal(t, kubeCfg.CurrentContext, changedKubeconfig.CurrentContext, "KubeConf.CurrentContext")
			}

		})
	}
}
