package tests

import (
	"os"
	"strings"
)

func GetTestWorkingDir() string {
	// to prevent access denied errors in github actions
	workspace := os.Getenv("GITHUB_WORKSPACE")
	if workspace == "" {
		workspace = "../tests/temp/"
	} else if !strings.HasSuffix(workspace, string(os.PathSeparator)) {
		workspace = workspace + string(os.PathSeparator)
	}
	return workspace
}
