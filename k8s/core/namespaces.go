package core

import (
	"context"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// NamespaceOps is an interface to perform namespace operations
type NamespaceOps interface {
	// ListNamespaces returns all the namespaces
	// Deprecated: Use ListNamespacesUsingLabelSelector instead.
	ListNamespaces(labelSelector map[string]string) (*corev1.NamespaceList, error)
	// ListNamespacesV2 returns all the namespaces when the labelSlector is a String
	// Deprecated: Use ListNamespacesUsingLabelSelector instead.
	ListNamespacesV2(labelSelector string) (*corev1.NamespaceList, error)
	// ListNamespacesUsingLabelSelector returns all the namespaces when the labelSlector is of type metav1.LabelSelector
	ListNamespacesUsingLabelSelector(labelSelector metav1.LabelSelector) (*corev1.NamespaceList, error)
	// GetNamespace returns a namespace object for given name
	GetNamespace(name string) (*corev1.Namespace, error)
	// CreateNamespace creates a namespace with given name and metadata
	CreateNamespace(*corev1.Namespace) (*corev1.Namespace, error)
	// UpdateNamespace update a namespace with given metadata
	UpdateNamespace(*corev1.Namespace) (*corev1.Namespace, error)
	// DeleteNamespace deletes a namespace with given name
	DeleteNamespace(name string) error
}

// ListNamespaces returns all the namespaces
// Deprecated: Use ListNamespacesUsingLabelSelector instead.
func (c *Client) ListNamespaces(labelSelector map[string]string) (*corev1.NamespaceList, error) {
	if err := c.initClient(); err != nil {
		return nil, err
	}

	return c.kubernetes.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{
		LabelSelector: mapToCSV(labelSelector),
	})
}

// ListNamespacesV2 returns all the namespaces when the labelSlector is a String
// Deprecated: Use ListNamespacesUsingLabelSelector instead.
func (c *Client) ListNamespacesV2(labelSelector string) (*corev1.NamespaceList, error) {
	if err := c.initClient(); err != nil {
		return nil, err
	}

	return c.kubernetes.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{
		LabelSelector: labelSelector,
	})
}

// ListNamespacesUsingLabelSelector list namespaces using the provided LabelSelector.
func (c *Client) ListNamespacesUsingLabelSelector(labelSelector metav1.LabelSelector) (*corev1.NamespaceList, error) {
	if err := c.initClient(); err != nil {
		return nil, err
	}

	namespaceLabelSelector, err := metav1.LabelSelectorAsSelector(&labelSelector)
	if err != nil {
		return nil, err
	}

	// List namespaces based on the label selector
	return c.kubernetes.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{
		LabelSelector: namespaceLabelSelector.String(),
	})
}

// GetNamespace returns a namespace object for given name
func (c *Client) GetNamespace(name string) (*corev1.Namespace, error) {
	if err := c.initClient(); err != nil {
		return nil, err
	}

	return c.kubernetes.CoreV1().Namespaces().Get(context.TODO(), name, metav1.GetOptions{})
}

// CreateNamespace creates a namespace with given name and metadata
func (c *Client) CreateNamespace(namespace *corev1.Namespace) (*corev1.Namespace, error) {
	if err := c.initClient(); err != nil {
		return nil, err
	}

	return c.kubernetes.CoreV1().Namespaces().Create(context.TODO(), namespace, metav1.CreateOptions{})
}

// DeleteNamespace deletes a namespace with given name
func (c *Client) DeleteNamespace(name string) error {
	if err := c.initClient(); err != nil {
		return err
	}

	return c.kubernetes.CoreV1().Namespaces().Delete(context.TODO(), name, metav1.DeleteOptions{})
}

// UpdateNamespace updates a namespace with given metadata
func (c *Client) UpdateNamespace(namespace *corev1.Namespace) (*corev1.Namespace, error) {
	if err := c.initClient(); err != nil {
		return nil, err
	}

	return c.kubernetes.CoreV1().Namespaces().Update(context.TODO(), namespace, metav1.UpdateOptions{})
}
