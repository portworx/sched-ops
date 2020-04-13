package networking

import (
	"time"

	v1beta1 "k8s.io/api/networking/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Ops is an interface to perform kubernetes related operations on the crd resources.
type IngressOps interface {
	// CreateIngress creates the given ingress
	CreateIngress(ingress *v1beta1.Ingress) (*v1beta1.Ingress, error)
	// UpdateIngress creates the given ingress
	UpdateIngress(ingress *v1beta1.Ingress) (*v1beta1.Ingress, error)
	// GetIngress returns the ingress given name and namespace
	GetIngress(name, namespace string) (*v1beta1.Ingress, error)
	// DeleteIngress deletes the given ingress
	DeleteIngress(name, namespace string) error
	// ValidateIngress validates the given ingress
	ValidateIngress(ingress *v1beta1.Ingress, timeout, retryInterval time.Duration) error
}

var NamespaceDefault = "default"

func (c *Client) CreateIngress(ingress *v1beta1.Ingress) (*v1beta1.Ingress, error) {
	if err := c.initClient(); err != nil {
		return nil, err
	}

	ns := ingress.Namespace
	if len(ns) == 0 {
		ns = NamespaceDefault
	}

	return c.networking.Ingresses(ns).Create(ingress)
}

func (c *Client) UpdateIngress(ingress *v1beta1.Ingress) (*v1beta1.Ingress, error) {
	if err := c.initClient(); err != nil {
		return nil, err
	}

	return c.networking.Ingresses(ingress.Namespace).Update(ingress)
}

func (c *Client) GetIngress(name, namespace string) (*v1beta1.Ingress, error) {
	if err := c.initClient(); err != nil {
		return nil, err
	}

	return c.networking.Ingresses(namespace).Get(name, metav1.GetOptions{})
}

func (c *Client) DeleteIngress(name, namespace string) error {
	if err := c.initClient(); err != nil {
		return err
	}

	return c.networking.Ingresses(namespace).Delete(name, &metav1.DeleteOptions{})
}

func (c *Client) ValidateIngress(ingress *v1beta1.Ingress, timeout, retryInterval time.Duration) error {
	if err := c.initClient(); err != nil {
		return err
	}

	result, err := c.networking.Ingresses(ingress.Namespace).Get(ingress.Name, metav1.GetOptions{})
	if result == nil {
		return err
	}
	return nil

}
