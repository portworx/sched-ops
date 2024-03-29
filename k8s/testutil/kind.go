package testutil

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"sync"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

const (
	kindKubeConfigPath      = "/tmp/kindutil.conf"
	envPreserveKindClusters = "PRESERVE_KIND_CLUSTERS"
)

// KindUtil interface is a wrapper around the kind utility
type KindUtil interface {
	ClusterExists(name string) (bool, error)
	CreateCluster(name string) error
	DestroyCluster(name string) error
	DestroyAllClusters() error
	GetClusterRestConfig(name string) (*rest.Config, error)
}

type kindUtil struct {
	kindPath string
}

var once sync.Once
var kind *kindUtil

// NewKindUtil initializes and returns a singleton kind util object
func NewKindUtil() KindUtil {
	once.Do(func() {
		kindPath, err := exec.LookPath("kind")
		if err != nil {
			logrus.Panicf("failed to locate kind: %v", err)
		}
		kind = &kindUtil{
			kindPath: kindPath,
		}
	})
	return kind
}

func (k *kindUtil) runCmd(cmd *exec.Cmd) (string, error) {
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("failed to run cmd %s: %v: %w", cmd.String(), string(out), err)
	}
	return string(out), nil
}

// CreateCluster creates a kind k8s cluster with the specified name
func (k *kindUtil) CreateCluster(name string) error {
	cmd := exec.Command(k.kindPath, "create", "cluster", "--name", name, "--wait=5m", "--kubeconfig", kindKubeConfigPath)
	_, err := k.runCmd(cmd)
	return err
}

// DestroyCluster destroys a kind k8s cluster with the specified name
func (k *kindUtil) DestroyCluster(name string) error {
	cmd := exec.Command(k.kindPath, "delete", "cluster", "--name", name, "--kubeconfig", kindKubeConfigPath)
	_, err := k.runCmd(cmd)
	return err
}

// DestroyAllClusters destroys all kind k8s clusters
func (k *kindUtil) DestroyAllClusters() error {
	cmd := exec.Command(k.kindPath, "delete", "clusters", "--all", "--kubeconfig", kindKubeConfigPath)
	_, err := k.runCmd(cmd)
	return err
}

// GetClusterRestConfig returns rest config for the specified kind cluster
func (k *kindUtil) GetClusterRestConfig(name string) (*rest.Config, error) {
	cmd := exec.Command(k.kindPath, "get", "kubeconfig", "--name", name)
	out, err := k.runCmd(cmd)
	if err != nil {
		return nil, fmt.Errorf("failed to get kubeconfig from kind: %v: %w", string(out), err)
	}
	clientCfg, err := clientcmd.NewClientConfigFromBytes([]byte(out))
	if err != nil {
		return nil, fmt.Errorf("failed to get clientCfg: %w", err)
	}
	restCfg, err := clientCfg.ClientConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to convert client config to rest config: %w", err)
	}
	return restCfg, nil
}

// ClusterExists checks if the specified kind cluster exists
func (k *kindUtil) ClusterExists(name string) (bool, error) {
	cmd := exec.Command(k.kindPath, "get", "clusters")
	out, err := k.runCmd(cmd)
	if err != nil {
		return false, fmt.Errorf("failed to list kind clusters: %w", err)
	}
	clusters := strings.Split(out, "\n")
	for _, cluster := range clusters {
		if cluster == name {
			return true, nil
		}
	}
	return false, nil
}

// SetUpTestCluster is a utility method to setup a kind cluster for testing
func SetUpTestCluster(t *testing.T, clusterName string) *rest.Config {
	// create a target cluster using kind if one does not exist already
	kind := NewKindUtil()
	exists, err := kind.ClusterExists(clusterName)
	require.Nil(t, err)
	if !exists {
		t.Logf("creating kind cluster %s", clusterName)
		err = kind.CreateCluster(clusterName)
		require.Nil(t, err)
		t.Cleanup(func() {
			DestroyTestCluster(t, clusterName)
		})
	}

	restCfg, err := kind.GetClusterRestConfig(clusterName)
	require.Nil(t, err)

	return restCfg
}

// DestroyTestCluster is a utility method to destroy a kind cluster from test
func DestroyTestCluster(t *testing.T, clusterName string) {
	var preserveKindClusters bool = true
	var err error

	val, defined := os.LookupEnv(envPreserveKindClusters)
	if defined {
		preserveKindClusters, err = strconv.ParseBool(val)
		require.Nil(t, err)
	}
	if preserveKindClusters {
		t.Logf("preserving kind cluster %s; use %s=false to change", clusterName, envPreserveKindClusters)
		return
	}
	kind := NewKindUtil()
	exists, err := kind.ClusterExists(clusterName)
	require.Nil(t, err)
	if exists {
		t.Logf("%s is set to %s, destroying kind cluster %s", envPreserveKindClusters, val, clusterName)
		err = kind.DestroyCluster(clusterName)
		require.Nil(t, err)
	}
}
