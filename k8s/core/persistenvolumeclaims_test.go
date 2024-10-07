package core

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestGetPersistentVolumeClaimsV2(t *testing.T) {
	// Create a mock Client
	client := MockClient()
	SetInstance(client)

	// Add PVC to the fake clientset
	_, err := client.kubernetes.CoreV1().PersistentVolumeClaims("test-namespace-1").Create(context.TODO(), &corev1.PersistentVolumeClaim{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-pvc-1",
			Namespace: "test-namespace-1",
			Labels:    map[string]string{"app": "test-app", "env": "test"},
		},
	}, metav1.CreateOptions{})
	require.NoError(t, err)

	_, err = client.kubernetes.CoreV1().PersistentVolumeClaims("test-namespace-1").Create(context.TODO(), &corev1.PersistentVolumeClaim{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-pvc-2",
			Namespace: "test-namespace-1",
			Labels:    map[string]string{"app": "test-app", "foo": "bar"},
		},
	}, metav1.CreateOptions{})
	require.NoError(t, err)

	// Test 1
	labelSelector := metav1.LabelSelector{
		MatchLabels: map[string]string{
			"app": "test-app",
		},
	}
	pvcList, err := client.GetPersistentVolumeClaimsUsingLabelSelector("test-namespace-1", labelSelector)
	require.NoError(t, err)
	assert.Len(t, pvcList.Items, 2)

	// Test 2
	labelSelector = metav1.LabelSelector{
		MatchLabels: map[string]string{
			"app": "test-app",
		},
		MatchExpressions: []metav1.LabelSelectorRequirement{
			{
				Key:      "env",
				Operator: metav1.LabelSelectorOpIn,
				Values:   []string{"test"},
			},
		},
	}
	pvcList, err = client.GetPersistentVolumeClaimsUsingLabelSelector("test-namespace-1", labelSelector)
	require.NoError(t, err)
	assert.Len(t, pvcList.Items, 1)
	assert.Equal(t, "test-pvc-1", pvcList.Items[0].Name)

	// Test 3
	labelSelector = metav1.LabelSelector{
		MatchExpressions: []metav1.LabelSelectorRequirement{
			{
				Key:      "foo",
				Operator: metav1.LabelSelectorOpDoesNotExist,
			},
		},
	}
	pvcList, err = client.GetPersistentVolumeClaimsUsingLabelSelector("test-namespace-1", labelSelector)
	require.NoError(t, err)
	assert.Len(t, pvcList.Items, 1)
	assert.Equal(t, "test-pvc-1", pvcList.Items[0].Name)
}
