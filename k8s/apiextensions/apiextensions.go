package apiextensions

import (
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/portworx/sched-ops/task"
	"github.com/sirupsen/logrus"
	apiextensionsclient "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/pkg/runtime"
	"k8s.io/client-go/pkg/watch"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

var (
	instance Ops
	once     sync.Once

	deleteForegroundPolicy = metav1.DeletePropagationForeground
)

// Ops is an interface to perform kubernetes related operations on the crd resources.
type Ops interface {
	CRDOps

	// SetConfig sets the config and resets the client.
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

// New builds a new apiextensions client.
func New(client apiextensionsclient.Interface) *Client {
	return &Client{
		extension: client,
	}
}

// NewForConfig builds a new apiextensions client for the given config.
func NewForConfig(c *rest.Config) (*Client, error) {
	client, err := apiextensionsclient.NewForConfig(c)
	if err != nil {
		return nil, err
	}

	return &Client{
		extension: client,
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

// Client provides a wrapper for kubernetes extension interface.
type Client struct {
	config    *rest.Config
	extension apiextensionsclient.Interface
}

// SetConfig sets the config and resets the client.
func (c *Client) SetConfig(cfg *rest.Config) {
	c.config = cfg
	c.extension = nil
}

// initClient the k8s client if uninitialized
func (c *Client) initClient() error {
	if c.extension != nil {
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

	c.extension, err = apiextensionsclient.NewForConfig(c.config)
	if err != nil {
		return err
	}

	return nil
}

// WatchFunc is a callback provided to the Watch functions
// which is invoked when the given object is changed.
type WatchFunc func(object runtime.Object) error

// handleWatch is internal function that handles the watch.  On channel shutdown (ie. stop watch),
// it'll attempt to reestablish its watch function.
func (c *Client) handleWatch(
	watchInterface watch.Interface,
	object runtime.Object,
	namespace string,
	fn WatchFunc,
	listOptions metav1.ListOptions) {
	defer watchInterface.Stop()
	for {
		select {
		case event, more := <-watchInterface.ResultChan():
			if !more {
				logrus.Debug("Kubernetes watch closed (attempting to re-establish)")

				t := func() (interface{}, bool, error) {
					var err error
					if _, ok := object.(*apiextensionsv1beta1.CustomResourceDefinition); ok {
						err = c.WatchCRDs(fn, listOptions)
					} else {
						return "", false, fmt.Errorf("unsupported object: %v given to handle watch", object)
					}

					return "", true, err
				}

				if _, err := task.DoRetryWithTimeout(t, 10*time.Minute, 10*time.Second); err != nil {
					logrus.WithError(err).Error("Could not re-establish the watch")
				} else {
					logrus.Debug("watch re-established")
				}
				return
			}

			fn(event.Object)
		}
	}
}
