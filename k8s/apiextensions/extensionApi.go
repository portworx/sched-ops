package apiextensions

import (
	"k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
)

// ExtensionOps is an interface to return api extension client interface
type ExtensionOps interface {
	// GetExtensionClient will return the extension client
	GetExtensionClient() (clientset.Interface, error)
}

// GetExtensionClient will return the extension client
func (c *Client) GetExtensionClient() (clientset.Interface, error) {
	if err := c.initClient(); err != nil {
		return nil, err
	}
	return c.extension, nil
}

