package discovery

import (
	"fmt"
	"os"
	"sync"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/version"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

var (
	instance Ops
	once     sync.Once

	deleteForegroundPolicy = metav1.DeletePropagationForeground
)

// Ops is an interface to the discovery client wrapper.
type Ops interface {
	Interface

	// SetConfig sets the config and resets the client
	SetConfig(config *rest.Config)
}

// Interface is an interface to perform kubernetes related operations on the discovery interface.
type Interface interface {
	// GetVersion gets the version from the kubernetes cluster
	GetVersion() (*version.Info, error)
}

// Instance returns a singleton instance of the client.
func Instance() Ops {
	once.Do(func() {
		if instance == nil {
			instance = &Client{}
		}
	})
	return instance
}

// SetInstance replaces the instance with the provided one. Should be used only for testing purposes.
func SetInstance(i Ops) {
	instance = i
}

// New builds a new client.
func New(c discovery.DiscoveryInterface) *Client {
	return &Client{
		discovery: c,
	}
}

// NewForConfig builds a new client for the given config.
func NewForConfig(c *rest.Config) (*Client, error) {
	discovery, err := discovery.NewDiscoveryClientForConfig(c)
	if err != nil {
		return nil, err
	}

	return &Client{
		discovery: discovery,
	}, nil
}

// Client is a wrapper for kubernetes discovery client.
type Client struct {
	config    *rest.Config
	discovery discovery.DiscoveryInterface
}

// SetConfig sets the config and resets the client.
func (c *Client) SetConfig(cfg *rest.Config) {
	c.config = cfg
	c.discovery = nil
}

func (c *Client) GetVersion() (*version.Info, error) {
	if err := c.initClient(); err != nil {
		return nil, err
	}

	return c.discovery.ServerVersion()
}

// initClient the k8s client if uninitialized
func (c *Client) initClient() error {
	if c.discovery != nil {
		return nil
	}

	return c.setClient()
}

// setClient instantiates a client.
func (c *Client) setClient() error {
	var err error

	if c.config != nil {
		err = c.loadClient()
	} else {
		kubeconfig := os.Getenv("KUBECONFIG")
		if len(kubeconfig) > 0 {
			err = c.loadClientFromKubeconfig(kubeconfig)
		} else {
			err = c.loadClientFromServiceAccount()
		}

	}

	return err
}

// loadClientFromServiceAccount loads a k8s client from a ServiceAccount specified in the pod running px
func (c *Client) loadClientFromServiceAccount() error {
	config, err := rest.InClusterConfig()
	if err != nil {
		return err
	}

	c.config = config
	return c.loadClient()
}

func (c *Client) loadClientFromKubeconfig(kubeconfig string) error {
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		return err
	}

	c.config = config
	return c.loadClient()
}

func (c *Client) loadClient() error {
	if c.config == nil {
		return fmt.Errorf("rest config is not provided")
	}

	var err error

	c.discovery, err = discovery.NewDiscoveryClientForConfig(c.config)
	if err != nil {
		return err
	}

	return nil
}
