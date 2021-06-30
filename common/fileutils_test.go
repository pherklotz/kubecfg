package common

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReadKubeConfigYaml_success(t *testing.T) {
	config, err := ReadKubeConfigYaml("../tests/resources/test.yaml")
	assert.Nil(t, err)
	assert.NotNil(t, config)
	assert.Equal(t, 2, len(config.Clusters))
	assert.Equal(t, 2, len(config.AuthInfos))
	assert.Equal(t, 2, len(config.Contexts))
}

func TestReadKubeConfigYaml_noFile(t *testing.T) {
	_, err := ReadKubeConfigYaml("../tests/resources/no_file.yaml")
	assert.NotNil(t, err)
}

func TestReadKubeConfigYaml_invalidFile(t *testing.T) {
	_, err := ReadKubeConfigYaml("../tests/resources/invalid.yaml")
	assert.NotNil(t, err)
}
