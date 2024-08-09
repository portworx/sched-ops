package configmap

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	coreops "github.com/portworx/sched-ops/k8s/core"
	fakek8sclient "k8s.io/client-go/kubernetes/fake"
)

func TestGetConfigMap(t *testing.T) {
	fakeClient := fakek8sclient.NewSimpleClientset()
	coreops.SetInstance(coreops.New(fakeClient))

	configData := map[string]string{
		"key1": "val1",
	}
	cm, err := New("px-configmaps-test", configData, testLockTimeout, 5, 0, 0)
	require.NoError(t, err, "Unexpected error in creating configmap")

	resultMap, err := cm.Get()
	require.NoError(t, err, "Unexpected error in getting configmap")
	require.Contains(t, resultMap, "key1")
	fmt.Println(resultMap)
}

func TestDeleteConfigMap(t *testing.T) {
	fakeClient := fakek8sclient.NewSimpleClientset()
	coreops.SetInstance(coreops.New(fakeClient))

	configData := map[string]string{
		"key1": "val1",
	}

	cm, err := New("px-configmaps-test", configData, testLockTimeout, 5, 0, 0)
	require.NoError(t, err, "Unexpected error in creating configmap")

	err = cm.Delete()
	require.NoError(t, err, "Unexpected error in delete")
}

func TestIncrementGeneration(t *testing.T) {
	fakeClient := fakek8sclient.NewSimpleClientset()
	coreops.SetInstance(coreops.New(fakeClient))

	configData := map[string]string{
		"key1": "1",
	}
	cmIntf, err := New("px-configmaps-test", configData, testLockTimeout, 5, 0, 0)
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
