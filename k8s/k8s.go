package k8s

import (
	"sync"

	prometheusclient "github.com/coreos/prometheus-operator/pkg/client/versioned"
	autopilotclientset "github.com/libopenstorage/autopilot-api/pkg/client/clientset/versioned"
	ostclientset "github.com/libopenstorage/operator/pkg/client/clientset/versioned"
	storkclientset "github.com/libopenstorage/stork/pkg/client/clientset/versioned"
	ocp_clientset "github.com/openshift/client-go/apps/clientset/versioned"
	ocp_security_clientset "github.com/openshift/client-go/security/clientset/versioned"
	"github.com/portworx/sched-ops/k8s/admissionregistration"
	"github.com/portworx/sched-ops/k8s/apiextensions"
	"github.com/portworx/sched-ops/k8s/apps"
	"github.com/portworx/sched-ops/k8s/autopilot"
	"github.com/portworx/sched-ops/k8s/batch"
	"github.com/portworx/sched-ops/k8s/core"
	"github.com/portworx/sched-ops/k8s/dynamic"
	"github.com/portworx/sched-ops/k8s/externalstorage"
	"github.com/portworx/sched-ops/k8s/openshift"
	"github.com/portworx/sched-ops/k8s/operator"
	"github.com/portworx/sched-ops/k8s/prometheus"
	"github.com/portworx/sched-ops/k8s/rbac"
	"github.com/portworx/sched-ops/k8s/storage"
	"github.com/portworx/sched-ops/k8s/stork"
	"github.com/portworx/sched-ops/k8s/talisman"
	talismanclientset "github.com/portworx/talisman/pkg/client/clientset/versioned"
	apiextensionsclient "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	dynamicclient "k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

// Ops is an interface to perform any kubernetes related operations
type Ops interface {
	admissionregistration.Interface
	apiextensions.Interface
	apps.Interface
	autopilot.Interface
	batch.Interface
	core.Interface
	dynamic.Interface
	externalstorage.Interface
	openshift.Interface
	operator.Interface
	prometheus.Interface
	rbac.Interface
	storage.Interface
	stork.Interface
	talisman.Interface
	ClientSetter
}

type AdmissionOps interface {
	admissionregistration.Interface
}

type ApiextensionsOps interface {
	apiextensions.Interface
}

type AppsOps interface {
	apps.Interface
}

type AutopilotOps interface {
	autopilot.Interface
}

type BatchOps interface {
	batch.Interface
}

type CoreOps interface {
	core.Interface
}

type DynamicOps interface {
	dynamic.Interface
}

type ExternalstorageOps interface {
	externalstorage.Interface
}

type OpenshiftOps interface {
	openshift.Interface
}

type OperatorOps interface {
	operator.Interface
}

type PrometheusOps interface {
	prometheus.Interface
}

type RbacOps interface {
	rbac.Interface
}

type StorageOps interface {
	storage.Interface
}

type StorkOps interface {
	stork.Interface
}

type TalismanOps interface {
	talisman.Interface
}

var (
	instance Ops
	once     sync.Once
)

type k8sOps struct {
	AdmissionOps
	ApiextensionsOps
	AppsOps
	AutopilotOps
	BatchOps
	CoreOps
	DynamicOps
	ExternalstorageOps
	OpenshiftOps
	OperatorOps
	PrometheusOps
	RbacOps
	StorageOps
	StorkOps
	TalismanOps

	client             kubernetes.Interface
	snapClient         rest.Interface
	storkClient        storkclientset.Interface
	ostClient          ostclientset.Interface
	talismanClient     talismanclientset.Interface
	autopilotClient    autopilotclientset.Interface
	apiExtensionClient apiextensionsclient.Interface
	dynamicInterface   dynamicclient.Interface
	ocpClient          ocp_clientset.Interface
	ocpSecurityClient  ocp_security_clientset.Interface
	prometheusClient   prometheusclient.Interface
}

// Instance returns a singleton instance of k8sOps type
func Instance() Ops {
	once.Do(func() {
		instance = &k8sOps{
			AdmissionOps:       &admissionregistration.Client{},
			ApiextensionsOps:   &apiextensions.Client{},
			AppsOps:            &apps.Client{},
			AutopilotOps:       &autopilot.Client{},
			BatchOps:           &batch.Client{},
			CoreOps:            &core.Client{},
			DynamicOps:         &dynamic.Client{},
			ExternalstorageOps: &externalstorage.Client{},
			OpenshiftOps:       &openshift.Client{},
			OperatorOps:        &operator.Client{},
			PrometheusOps:      &prometheus.Client{},
			RbacOps:            &rbac.Client{},
			StorageOps:         &storage.Client{},
			StorkOps:           &stork.Client{},
			TalismanOps:        &talisman.Client{},
		}
	})
	return instance
}

// NewInstanceFromClients returns new instance of k8sOps by using given
// clients
func NewInstanceFromClients(
	kubernetesClient kubernetes.Interface,
	snapClient rest.Interface,
	storkClient storkclientset.Interface,
	apiExtensionClient apiextensionsclient.Interface,
	dynamicClient dynamicclient.Interface,
	ocpClient ocp_clientset.Interface,
	ocpSecurityClient ocp_security_clientset.Interface,
	autopilotClient autopilotclientset.Interface,
) Ops {
	c := &k8sOps{
		client:             kubernetesClient,
		snapClient:         snapClient,
		storkClient:        storkClient,
		apiExtensionClient: apiExtensionClient,
		dynamicInterface:   dynamicClient,
		ocpClient:          ocpClient,
		ocpSecurityClient:  ocpSecurityClient,
		autopilotClient:    autopilotClient,
	}
	c.setClients()
	return c
}

// NewInstanceFromRestConfig returns new instance of k8sOps by using given
// k8s rest client config
func NewInstanceFromRestConfig(config *rest.Config) (Ops, error) {
	k := &k8sOps{}
	if err := k.loadClientFor(config); err != nil {
		return nil, err
	}
	return k, nil
}
