package common

import (
	"testing"

	"github.com/stretchr/testify/assert"
	k8s "k8s.io/client-go/tools/clientcmd/api/v1"
)

func TestGetContextByName(t *testing.T) {
	contextName := "test-ctx"
	expectedCtx := k8s.NamedContext{
		Name: contextName,
		Context: k8s.Context{
			Cluster:  "test cluster",
			AuthInfo: "test auth info",
		},
	}

	config := k8s.Config{
		Contexts: []k8s.NamedContext{
			{
				Name: "test ctx",
				Context: k8s.Context{
					Cluster:  "test cluster",
					AuthInfo: "test auth info",
				},
			},
			expectedCtx,
			{
				Name: "contextName",
				Context: k8s.Context{
					Cluster:  "another test cluster",
					AuthInfo: "test auth info",
				},
			},
		},
	}

	actualCtx, err := GetContextByName(&config, &contextName)

	assert.Nil(t, err)
	assert.ObjectsAreEqual(expectedCtx, actualCtx)
}

func TestGetContextByName_failed(t *testing.T) {
	contextName := "test-ctx"

	config := k8s.Config{
		Contexts: []k8s.NamedContext{
			{
				Name: "test-context",
				Context: k8s.Context{
					Cluster:  "test cluster",
					AuthInfo: "test auth info",
				},
			},
			{
				Name: "another-test",
				Context: k8s.Context{
					Cluster:  "another test cluster",
					AuthInfo: "test auth info",
				},
			},
		},
	}

	_, err := GetContextByName(&config, &contextName)

	assert.NotNil(t, err)
	assert.EqualError(t, err, "Can not find context with name '"+contextName+"'")
}

func TestGetClusterByName(t *testing.T) {
	name := "test-cluster"
	expected := k8s.NamedCluster{
		Name: name,
		Cluster: k8s.Cluster{
			Server: "test server",
		},
	}

	config := k8s.Config{
		Clusters: []k8s.NamedCluster{
			{
				Name: "test ctx",
				Cluster: k8s.Cluster{
					Server: "another test server",
				},
			},
			expected,
			{
				Name: "clusterName",
				Cluster: k8s.Cluster{
					Server: "another test server asd",
				},
			},
		},
	}

	actualCtx, err := GetClusterByName(&config, &name)

	assert.Nil(t, err)
	assert.ObjectsAreEqual(expected, actualCtx)
}

func TestGetClusterByName_failed(t *testing.T) {
	name := "test-ctx"

	config := k8s.Config{
		Clusters: []k8s.NamedCluster{
			{
				Name: "test ctx",
				Cluster: k8s.Cluster{
					Server: "another test server",
				},
			},
			{
				Name: "clusterName",
				Cluster: k8s.Cluster{
					Server: "another test server asd",
				},
			},
		},
	}

	_, err := GetClusterByName(&config, &name)

	assert.NotNil(t, err)
	assert.EqualError(t, err, "Can not find cluster with name '"+name+"'")
}

func TestGetUserByName(t *testing.T) {
	name := "test-auth-info"
	expected := k8s.NamedAuthInfo{
		Name: name,
		AuthInfo: k8s.AuthInfo{
			Username: "test user",
		},
	}

	config := k8s.Config{
		AuthInfos: []k8s.NamedAuthInfo{
			{
				Name: "test auth info",
				AuthInfo: k8s.AuthInfo{
					Username: "test user",
				},
			},
			expected,
			{
				Name: "auth info",
				AuthInfo: k8s.AuthInfo{
					Username: "test user",
				},
			},
		},
	}

	actualCtx, err := GetUserByName(&config, &name)

	assert.Nil(t, err)
	assert.ObjectsAreEqual(expected, actualCtx)
}

func TestGetUserByName_failed(t *testing.T) {
	name := "test-auth-info"

	config := k8s.Config{
		AuthInfos: []k8s.NamedAuthInfo{
			{
				Name: "test auth info",
				AuthInfo: k8s.AuthInfo{
					Username: "test user",
				},
			},
			{
				Name: "auth info",
				AuthInfo: k8s.AuthInfo{
					Username: "test user",
				},
			},
		},
	}

	_, err := GetUserByName(&config, &name)

	assert.NotNil(t, err)
	assert.EqualError(t, err, "Can not find user with name '"+name+"'")
}
