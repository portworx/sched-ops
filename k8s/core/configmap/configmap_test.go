package configmap

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	coreops "github.com/portworx/sched-ops/k8s/core"
	"github.com/portworx/sched-ops/k8s/testutil"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestGetConfigMap(t *testing.T) {
	setUpConfigMapTestCluster(t)

	configData := map[string]string{
		"key1": "val1",
	}
	cm, err := New("px-configmaps-get-test", configData, testLockTimeout, 5, 0, 0)
	require.NoError(t, err, "Unexpected error in creating configmap")

	resultMap, err := cm.Get()
	require.NoError(t, err, "Unexpected error in getting configmap")
	require.Contains(t, resultMap, "key1")
	fmt.Println(resultMap)
}

func TestDeleteConfigMap(t *testing.T) {
	setUpConfigMapTestCluster(t)

	configData := map[string]string{
		"key1": "val1",
	}

	cm, err := New("px-configmaps-delete-test", configData, testLockTimeout, 5, 0, 0)
	require.NoError(t, err, "Unexpected error in creating configmap")

	err = cm.Delete()
	require.NoError(t, err, "Unexpected error in delete")
}

func TestIncrementGeneration(t *testing.T) {
	setUpConfigMapTestCluster(t)

	configData := map[string]string{
		"key1": "1",
	}
	cmIntf, err := New("px-configmaps-increment-generation-test", configData, testLockTimeout, 5, 0, 0)
	require.NoError(t, err, "Unexpected error in creating configmap")

	cm := cmIntf.(*configMap)

	rawCM, err := coreops.Instance().GetConfigMap(cm.name, k8sSystemNamespace)
	require.NoError(t, err, "Unexpected error in getting raw configmap")

	require.Equal(t, "", rawCM.Data[pxGenerationKey])
	newGen := cm.incrementGeneration(rawCM)
	require.Equal(t, "1", rawCM.Data[pxGenerationKey])
	require.Equal(t, uint64(1), newGen)

	newGen = cm.incrementGeneration(rawCM)
	require.Equal(t, "2", rawCM.Data[pxGenerationKey])
	require.Equal(t, uint64(2), newGen)

	rawCM.Data[pxGenerationKey] = "123456789"
	newGen = cm.incrementGeneration(rawCM)
	require.Equal(t, "123456790", rawCM.Data[pxGenerationKey])
	require.Equal(t, uint64(123456790), newGen)

	rawCM.Data[pxGenerationKey] = "invalid"
	newGen = cm.incrementGeneration(rawCM)
	require.Equal(t, "1", rawCM.Data[pxGenerationKey])
	require.Equal(t, uint64(1), newGen)

	// max uint64
	rawCM.Data[pxGenerationKey] = "18446744073709551615"
	newGen = cm.incrementGeneration(rawCM)
	require.Equal(t, "1", rawCM.Data[pxGenerationKey])
	require.Equal(t, uint64(1), newGen)

	err = cm.Delete()
	require.NoError(t, err, "Unexpected error in delete")
}

func setUpConfigMapTestCluster(t *testing.T) {
	os.Setenv("KUBERNETES_OPS_QPS_RATE", "2000")
	os.Setenv("KUBERNETES_OPS_BURST_RATE", "4000")

	restCfg := testutil.SetUpTestCluster(t, "configmap-test-cluster")

	testClient, err := coreops.NewForConfig(restCfg)
	require.NoError(t, err)

	coreops.SetInstance(testClient)
	// delete all the test configmaps
	result, err := coreops.Instance().ListConfigMap(k8sSystemNamespace, metav1.ListOptions{})
	require.NoError(t, err)
	for _, cm := range result.Items {
		if strings.Contains(cm.Name, "px-configmaps") && strings.Contains(cm.Name, "test") {
			t.Logf("Deleting configmap: %s/%s", k8sSystemNamespace, cm.Name)
			err = coreops.Instance().DeleteConfigMap(cm.Name, k8sSystemNamespace)
			require.NoError(t, err)
		}
	}
}
