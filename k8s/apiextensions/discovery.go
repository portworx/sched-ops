package apiextensions

import (
	"k8s.io/client-go/discovery"
)

// ExtensionOps is an interface to return discovery client interface
type DiscoveryOps interface {
	// GetDiscoveryClient will return the extension client
	GetDiscoveryClient() (discovery.DiscoveryInterface, error)
}

// GetDiscoveryClient will return the discovery client interface
func (c *Client) GetDiscoveryClient() (discovery.DiscoveryInterface, error) {
	if err := c.initClient(); err != nil {
		return nil, err
	}
	return c.extension.Discovery(), nil
}

