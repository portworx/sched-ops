package rancher

import (
	"sync"

	"github.com/portworx/sched-ops/k8s/common"
	rancherv3 "github.com/rancher/rancher/pkg/apis/management.cattle.io/v3"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
)

var (
	instance Ops
	once     sync.Once
)

// Ops is an interface to Operator operations.
type Ops interface {
	ProjectOps

	// SetConfig sets the config and resets the client
	SetConfig(config *rest.Config)
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

// New builds a new rancher client.
func New(crClient client.Client) *Client {
	return &Client{
		crClient: crClient,
	}
}

// Client is a wrapper for the rancher client.
type Client struct {
	config   *rest.Config
	crClient client.Client
}

// SetConfig sets the config and resets the client
func (c *Client) SetConfig(cfg *rest.Config) {
	c.config = cfg
	c.crClient = nil
}

// initClient the k8s client if uninitialized
func (c *Client) initClient() error {
	if c.crClient != nil {
		return nil
	}

	return c.setClient()
}

// setClient instantiates a client.
func (c *Client) setClient() error {
	if c.config == nil {
		c.config = config.GetConfigOrDie()
	}

	err := common.SetRateLimiter(c.config)
	if err != nil {
		return err
	}

	s := scheme.Scheme
	if err := rancherv3.AddToScheme(s); err != nil {
		return err
	}

	c.crClient, err = client.New(c.config, client.Options{Scheme: s})
	if err != nil {
		return err
	}

	return nil
}
