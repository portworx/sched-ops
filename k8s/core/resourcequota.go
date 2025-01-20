package core

import (
	"context"

	"github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
)

// ResourceQuotaOps is an interface to perform k8s ResourceQuota operations
type ResourceQuotaOps interface {
	// GetResourceQuota gets the resource quota object for the given name and namespace
	GetResourceQuota(name string, namespace string) (*corev1.ResourceQuota, error)
	// CreateResourceQuota creates a new resource quota object if it does not already exist.
	CreateResourceQuota(resourceQuota *corev1.ResourceQuota) (*corev1.ResourceQuota, error)
	// DeleteResourceQuota deletes the given resource quota
	DeleteResourceQuota(name, namespace string) error
	// UpdateResourceQuota updates the given resource quota object
	UpdateResourceQuota(resourceQuota *corev1.ResourceQuota) (*corev1.ResourceQuota, error)
	// WatchResourceQuota sets up a watcher that listens for changes on the resource quota
	WatchResourceQuota(resourceQuota *corev1.ResourceQuota, fn WatchFunc) error
	// ListResourceQuota returns the list of ResourceQuotas
	ListResourceQuota(namespace string, filterOptions metav1.ListOptions) (*corev1.ResourceQuotaList, error)
}

// GetResourceQuota gets the resource quota object for the given name and namespace
func (c *Client) GetResourceQuota(name string, namespace string) (*corev1.ResourceQuota, error) {
	if err := c.initClient(); err != nil {
		return nil, err
	}

	return c.kubernetes.CoreV1().ResourceQuotas(namespace).Get(context.TODO(), name, metav1.GetOptions{})
}

// CreateResourceQuota creates a new resource quota object if it does not already exist.
func (c *Client) CreateResourceQuota(resourceQuota *corev1.ResourceQuota) (*corev1.ResourceQuota, error) {
	if err := c.initClient(); err != nil {
		return nil, err
	}

	ns := resourceQuota.Namespace
	if len(ns) == 0 {
		ns = corev1.NamespaceDefault
	}

	return c.kubernetes.CoreV1().ResourceQuotas(ns).Create(context.TODO(), resourceQuota, metav1.CreateOptions{})
}

// DeleteResourceQuota deletes the given resource quota
func (c *Client) DeleteResourceQuota(name, namespace string) error {
	if err := c.initClient(); err != nil {
		return err
	}

	if len(namespace) == 0 {
		namespace = corev1.NamespaceDefault
	}

	return c.kubernetes.CoreV1().ResourceQuotas(namespace).Delete(context.TODO(), name, metav1.DeleteOptions{
		PropagationPolicy: &deleteForegroundPolicy,
	})
}

// UpdateResourceQuota updates the given resource quota object
func (c *Client) UpdateResourceQuota(resourceQuota *corev1.ResourceQuota) (*corev1.ResourceQuota, error) {
	if err := c.initClient(); err != nil {
		return nil, err
	}

	ns := resourceQuota.Namespace
	if len(ns) == 0 {
		ns = corev1.NamespaceDefault
	}

	return c.kubernetes.CoreV1().ResourceQuotas(ns).Update(context.TODO(), resourceQuota, metav1.UpdateOptions{})
}

// WatchResourceQuota sets up a watcher that listens for changes on the resource quota
func (c *Client) WatchResourceQuota(resourceQuota *corev1.ResourceQuota, fn WatchFunc) error {
	if err := c.initClient(); err != nil {
		return err
	}

	listOptions := metav1.ListOptions{
		FieldSelector: fields.OneTermEqualSelector("metadata.name", resourceQuota.Name).String(),
		Watch:         true,
	}

	watchInterface, err := c.kubernetes.CoreV1().ResourceQuotas(resourceQuota.Namespace).Watch(context.TODO(), listOptions)
	if err != nil {
		logrus.WithError(err).Error("error invoking the watch api for resource quotas")
		return err
	}

	// fire off watch function
	go c.handleWatch(watchInterface, resourceQuota, "", fn, listOptions)
	return nil
}

// ListResourceQuota returns the list of ResourceQuotas
func (c *Client) ListResourceQuota(namespace string, filterOptions metav1.ListOptions) (*corev1.ResourceQuotaList, error) {
	if err := c.initClient(); err != nil {
		return nil, err
	}

	return c.kubernetes.CoreV1().ResourceQuotas(namespace).List(context.TODO(), filterOptions)
}
