package core

import (
	"context"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/events"
	"k8s.io/client-go/tools/record"
)

// MockClient returns a mock Client for testing
func MockClient() *Client {
	// Create a fake Kubernetes client
	fakeClientset := fake.NewSimpleClientset()

	// Create the mock client struct
	return &Client{
		config:                 &rest.Config{},                        // Mock rest.Config if needed
		kubernetes:             fakeClientset,                         // Use fake clientset as Kubernetes interface
		eventRecordersLegacy:   make(map[string]record.EventRecorder), // Mock legacy recorders
		eventBroadcasterLegacy: nil,                                   // No-op for legacy broadcaster
		eventRecordersNew:      make(map[string]events.EventRecorder), // Mock new recorders
		eventBroadcasterNew:    nil,                                   // No-op for new broadcaster
		eventRecordersLock:     sync.Mutex{},                          // Mock sync lock
	}
}

// Setup cluster ad write a test to list namespaces
func TestListNamespacesV3(t *testing.T) {
	// Create a mock Client
	client := MockClient()
	SetInstance(client)

	// Add namespace to the fake clientset
	_, err := client.kubernetes.CoreV1().Namespaces().Create(context.TODO(), &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name:   "test-namespace-1",
			Labels: map[string]string{"app": "test-app", "env": "test"},
		},
	}, metav1.CreateOptions{})
	require.NoError(t, err)

	_, err = client.kubernetes.CoreV1().Namespaces().Create(context.TODO(), &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name:   "test-namespace-2",
			Labels: map[string]string{"app": "test-app", "foo": "bar"},
		},
	}, metav1.CreateOptions{})
	require.NoError(t, err)

	// Test 1
	labelSelector := metav1.LabelSelector{
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
	namespaceList, err := client.ListNamespacesV3(labelSelector)
	require.NoError(t, err)
	assert.Len(t, namespaceList.Items, 1)
	assert.Equal(t, "test-namespace-1", namespaceList.Items[0].Name)

	// Test 2
	labelSelector = metav1.LabelSelector{
		MatchLabels: map[string]string{
			"app": "test-app",
		},
		MatchExpressions: []metav1.LabelSelectorRequirement{
			{
				Key:      "foo",
				Operator: metav1.LabelSelectorOpIn,
				Values:   []string{"bar"},
			},
		},
	}
	namespaceList, err = client.ListNamespacesV3(labelSelector)
	require.NoError(t, err)
	assert.Len(t, namespaceList.Items, 1)

	// Test 3
	labelSelector = metav1.LabelSelector{
		MatchExpressions: []metav1.LabelSelectorRequirement{
			{
				Key:      "foo",
				Operator: metav1.LabelSelectorOpDoesNotExist,
			},
		},
	}
	namespaceList, err = client.ListNamespacesV3(labelSelector)
	require.NoError(t, err)
	assert.Len(t, namespaceList.Items, 1)

	// Test 4
	labelSelector = metav1.LabelSelector{
		MatchExpressions: []metav1.LabelSelectorRequirement{
			{
				Key:      "foo",
				Operator: metav1.LabelSelectorOpExists,
			},
		},
	}
	namespaceList, err = client.ListNamespacesV3(labelSelector)
	require.NoError(t, err)
	assert.Len(t, namespaceList.Items, 1)

	// Test 5
	labelSelector = metav1.LabelSelector{
		MatchLabels: map[string]string{
			"app": "test-app",
		},
	}
	namespaceList, err = client.ListNamespacesV3(labelSelector)
	require.NoError(t, err)
	assert.Len(t, namespaceList.Items, 2)

	//Test 6
	labelSelector = metav1.LabelSelector{}
	namespaceList, err = client.ListNamespacesV3(labelSelector)
	require.NoError(t, err)
	assert.Len(t, namespaceList.Items, 2)

}
