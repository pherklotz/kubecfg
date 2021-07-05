package common

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	k8s "k8s.io/client-go/tools/clientcmd/api/v1"
)

const TEST_RESOURCE_DIR = "../tests/resources/"

func getTestWorkingDir() string {
	// to prevent access denied errors in github actions
	workspace := os.Getenv("GITHUB_WORKSPACE")
	if workspace == "" {
		workspace = "../tests/temp/"
	} else if !strings.HasSuffix(workspace, string(os.PathSeparator)) {
		workspace = workspace + string(os.PathSeparator)
	}
	return workspace
}

func TestReadKubeConfigYaml_success(t *testing.T) {
	config, err := ReadKubeConfigYaml(TEST_RESOURCE_DIR + "test.yaml")
	assert.Nil(t, err)
	assert.NotNil(t, config)
	assert.Equal(t, 2, len(config.Clusters))
	assert.Equal(t, 2, len(config.AuthInfos))
	assert.Equal(t, 2, len(config.Contexts))
}

func TestReadKubeConfigYaml_noFile(t *testing.T) {
	_, err := ReadKubeConfigYaml(TEST_RESOURCE_DIR + "no_file.yaml")
	assert.NotNil(t, err)
}

func TestReadKubeConfigYaml_invalidFile(t *testing.T) {
	_, err := ReadKubeConfigYaml(TEST_RESOURCE_DIR + "invalid.yaml")
	assert.NotNil(t, err)
}

func TestFileExists(t *testing.T) {
	var tests = []struct {
		filePath string
		expected bool
	}{
		{TEST_RESOURCE_DIR + "not_existing_file", false},
		{TEST_RESOURCE_DIR + "test.yaml", true},
	}
	for _, tt := range tests {
		testname := fmt.Sprintf("Expect '%t' on '%s' call of FileExists", tt.expected, tt.filePath)
		t.Run(testname, func(t *testing.T) {
			actual := FileExists(tt.filePath)
			assert.Equal(t, tt.expected, actual)
		})
	}
}

func TestIsRegularFile(t *testing.T) {
	var tests = []struct {
		filePath string
		expected bool
	}{
		{TEST_RESOURCE_DIR, false},
		{TEST_RESOURCE_DIR + "test.yaml", true},
	}
	for _, tt := range tests {
		testname := fmt.Sprintf("Expect '%t' on '%s' call of IsRegularFile", tt.expected, tt.filePath)
		t.Run(testname, func(t *testing.T) {
			actual := IsRegularFile(tt.filePath)
			assert.Equal(t, tt.expected, actual)
		})
	}
}

func TestCopyFile(t *testing.T) {
	source := TEST_RESOURCE_DIR + "test.yaml"
	targetDir := getTestWorkingDir()
	err := os.MkdirAll(targetDir, os.ModeDir)
	assert.Nil(t, err)

	target := targetDir + "test_copy.yaml"
	err = CopyFile(source, target)
	assert.Nil(t, err)
	defer os.Remove(target)
	targetInfo, err := os.Stat(target)
	assert.Nil(t, err)
	sourceInfo, err := os.Stat(source)
	assert.Nil(t, err)
	assert.Equal(t, sourceInfo.Size(), targetInfo.Size())
}

func TestWriteAndReadKubeConfigYaml(t *testing.T) {
	target := getTestWorkingDir() + "new_config.yaml"
	defer os.Remove(target)
	expectedConfig := k8s.Config{}
	WriteKubeConfigYaml(target, &expectedConfig)
	actualConfig, err := ReadKubeConfigYaml(target)
	assert.Nil(t, err)

	assert.ObjectsAreEqual(expectedConfig, actualConfig)
}
