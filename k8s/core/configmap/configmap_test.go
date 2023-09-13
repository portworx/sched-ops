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
	cm, err := New("px-configmaps-test", configData, lockTimeout, 5, 0, 0, "test-ns")
	fmt.Println("cm : ", cm)
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

	cm, err := New("px-configmaps-test", configData, lockTimeout, 5, 0, 0, "test-ns")
	require.NoError(t, err, "Unexpected error in creating configmap")

	err = cm.Delete()
	require.NoError(t, err, "Unexpected error in delete")

}

func TestPatchConfigMap(t *testing.T) {
	fakeClient := fakek8sclient.NewSimpleClientset()
	coreops.SetInstance(coreops.New(fakeClient))

	configData := map[string]string{
		"key1": "val1",
	}

	cm, err := New("px-configmaps-test", configData, lockTimeout, 5, 0, 0, "test-ns")
	require.NoError(t, err, "Unexpected error in creating configmap")

	dummyData := map[string]string{
		"key2": "val2",
	}

	err = cm.Patch(dummyData)
	require.NoError(t, err, "Unexpected error in Patch")
	resultMap, err := cm.Get()
	require.Contains(t, resultMap, "key1")
	require.Contains(t, resultMap, "key2")
	fmt.Println(resultMap)
}

func TestUpdateConfigMap(t *testing.T) {
	fakeClient := fakek8sclient.NewSimpleClientset()
	coreops.SetInstance(coreops.New(fakeClient))

	configData := map[string]string{
		"key1": "val1",
	}

	cm, err := New("px-configmaps-test", configData, lockTimeout, 5, 0, 0, "test-ns")
	require.NoError(t, err, "Unexpected error in creating configmap")

	dummyData := map[string]string{
		"key2": "val2",
	}

	err = cm.Update(dummyData)
	require.NoError(t, err, "Unexpected error in Update")
	resultMap, err := cm.Get()
	require.NotContains(t, resultMap, "key1")
	require.Contains(t, resultMap, "key2")
	fmt.Println(resultMap)
}
