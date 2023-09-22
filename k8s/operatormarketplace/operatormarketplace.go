package operatormarketplace

import (
	"fmt"
	"sync"

	ofv1 "github.com/operator-framework/api/pkg/operators/v1"
	ofv1alpha1 "github.com/operator-framework/api/pkg/operators/v1alpha1"
	"github.com/portworx/sched-ops/k8s/common"
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
	CatalogSourceOps
	ClusterServiceVersionOps
	OperatorGroupOps
	SubscriptionOps
	InstallPlanOps

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

// New builds a new operator client.
func New(crClient client.Client) *Client {
	return &Client{
		crClient: crClient,
	}
}

// NewForConfig builds a new client for the given config.
func NewForConfig(c *rest.Config) (*Client, error) {
	crClient, err := crClient.NewForConfig(c)
	if err != nil {
		return nil, err
	}

	return &Client{
		crClient: crClient,
	}, nil
}

// NewInstanceFromConfigFile returns new instance of client by using given
// config file
func NewInstanceFromConfigFile(config string) (Ops, error) {
	newInstance := &Client{}
	err := newInstance.loadClientFromKubeconfig(config)
	if err != nil {
		return nil, err
	}
	return newInstance, nil
}

// Client is a wrapper for the operator client.
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
	if err := ofv1alpha1.AddToScheme(s); err != nil {
		return err
	}
	if err := ofv1.AddToScheme(s); err != nil {
		return err
	}
	c.crClient, err = client.New(c.config, client.Options{Scheme: s})
	if err != nil {
		return err
	}

	return nil
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
	err = common.SetRateLimiter(c.config)
	if err != nil {
		return err
	}
	c.crClient, err = crClient.NewForConfig(c.config)
	if err != nil {
		return err
	}
	return nil
}
