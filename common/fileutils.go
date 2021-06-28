package common

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"

	"gopkg.in/yaml.v2"
	k8s "k8s.io/client-go/tools/clientcmd/api/v1"
)

// ReadKubeConfigYaml deserializes a kubeconfig yaml file into a KubeConfig object.
func ReadKubeConfigYaml(filePath string) (kc *k8s.Config, err error) {
	err = IsRegularFile(filePath)
	if err != nil {
		return
	}

	fileContent, err := ioutil.ReadFile(filePath)
	if err != nil {
		return
	}
	kc = &k8s.Config{}
	err = yaml.Unmarshal(fileContent, kc)
	return
}

// FileExists reports whether the named file or directory exists.
func FileExists(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

// IsRegularFile checks whether the path is a regular file or not.
func IsRegularFile(path string) error {
	sourceFileStat, err := os.Stat(path)
	if err != nil {
		return err
	}
	if !sourceFileStat.Mode().IsRegular() {
		return fmt.Errorf("'%s' is not a regular file", path)
	}
	return nil
}

// CopyFile copies a file
func CopyFile(src, target string) error {
	err := IsRegularFile(src)
	if err != nil {
		return err
	}

	source, err := os.Open(src)
	if err != nil {
		return err
	}
	defer source.Close()

	targetFile, err := os.Create(target)
	if err != nil {
		return err
	}
	defer targetFile.Close()
	_, err = io.Copy(targetFile, source)
	return err
}

// GetDefaultKubeconfigPath returns the path to the default user kubeconfig (e.g. ~/.kube/config)
func GetDefaultKubeconfigPath() (path string, err error) {
	usrHome, err := os.UserHomeDir()
	if err == nil {
		path = usrHome + string(os.PathSeparator) + ".kube" + string(os.PathSeparator) + "config"
	}
	return
}

// GetDefaultKubeconfig returns the default user kubeconfig (e.g. ~/.kube/config)
func GetDefaultKubeconfig() (*k8s.Config, error) {
	usrHome, err := os.UserHomeDir()
	if err != nil {
		return &k8s.Config{}, err
	}
	path := usrHome + string(os.PathSeparator) + ".kube" + string(os.PathSeparator) + "config"
	return ReadKubeConfigYaml(path)
}

// WriteKubeConfigYaml writes a KubeConfig into a file
func WriteKubeConfigYaml(target string, conf *k8s.Config) {
	targetYaml, err := yaml.Marshal(&conf)
	if err != nil {
		log.Fatalln("can not marshal kubeconfig yaml: ", err)
	}
	err = ioutil.WriteFile(target, targetYaml, 0600)
	if err != nil {
		log.Fatalf("Can not write target yaml '%s': %v\n", target, err)
	}
}
