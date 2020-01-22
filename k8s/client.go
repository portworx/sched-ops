package k8s

import (
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
	"github.com/portworx/sched-ops/k8s/discovery"
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
	"k8s.io/client-go/tools/clientcmd"
)

// ClientSetter is an interface to allow setting different clients on the Ops object
type ClientSetter interface {
	// SetConfig sets the config and resets the client
	SetConfig(config *rest.Config)
	// SetConfigFromPath sets the config from a kubeconfig file
	SetConfigFromPath(configPath string) error
	// SetClient set the k8s clients
	SetClient(
		kubernetes.Interface,
		rest.Interface,
		storkclientset.Interface,
		apiextensionsclient.Interface,
		dynamicclient.Interface,
		ocp_clientset.Interface,
		ocp_security_clientset.Interface,
		autopilotclientset.Interface,
	)
	// SetBaseClient sets the kubernetes clientset
	SetBaseClient(kubernetes.Interface)
	// SetSnapshotClient sets the snapshot clientset
	SetSnapshotClient(rest.Interface)
	// SetStorkClient sets the stork clientset
	SetStorkClient(storkclientset.Interface)
	// SetOpenstorageOperatorClient sets the openstorage operator clientset
	SetOpenstorageOperatorClient(ostclientset.Interface)
	// SetAPIExtensionsClient sets the api extensions clientset
	SetAPIExtensionsClient(apiextensionsclient.Interface)
	// SetDynamicClient sets the dynamic clientset
	SetDynamicClient(dynamicclient.Interface)
	// SetOpenshiftAppsClient sets the openshift apps clientset
	SetOpenshiftAppsClient(ocp_clientset.Interface)
	// SetOpenshiftSecurityClient sets the openshift security clientset
	SetOpenshiftSecurityClient(ocp_security_clientset.Interface)
	// SetTalismanClient sets the talisman clientset
	SetTalismanClient(talismanclientset.Interface)
	// SetAutopilotClient sets the autopilot clientset
	SetAutopilotClient(autopilotclientset.Interface)
	// SetPrometheusClient sets the prometheus clientset
	SetPrometheusClient(prometheusclient.Interface)
}

// SetConfig sets the config and resets the client
func (k *k8sOps) SetConfig(config *rest.Config) {
	admissionregistration.Instance().SetConfig(config)
	apiextensions.Instance().SetConfig(config)
	apps.Instance().SetConfig(config)
	autopilot.Instance().SetConfig(config)
	batch.Instance().SetConfig(config)
	core.Instance().SetConfig(config)
	discovery.Instance().SetConfig(config)
	dynamic.Instance().SetConfig(config)
	externalstorage.Instance().SetConfig(config)
	openshift.Instance().SetConfig(config)
	operator.Instance().SetConfig(config)
	prometheus.Instance().SetConfig(config)
	rbac.Instance().SetConfig(config)
	storage.Instance().SetConfig(config)
	stork.Instance().SetConfig(config)
	talisman.Instance().SetConfig(config)
}

// SetConfigFromPath takes the path to a kubeconfig file
// and then internally calls SetConfig to set it
func (k *k8sOps) SetConfigFromPath(configPath string) error {
	if configPath == "" {
		k.SetConfig(nil)
		return nil
	}
	config, err := clientcmd.BuildConfigFromFlags("", configPath)
	if err != nil {
		return err
	}

	k.SetConfig(config)
	return nil
}

// SetClient set the k8s clients
func (k *k8sOps) SetClient(
	client kubernetes.Interface,
	snapClient rest.Interface,
	storkClient storkclientset.Interface,
	apiExtensionClient apiextensionsclient.Interface,
	dynamicInterface dynamicclient.Interface,
	ocpClient ocp_clientset.Interface,
	ocpSecurityClient ocp_security_clientset.Interface,
	autopilotClient autopilotclientset.Interface,
) {
	k.client = client
	k.snapClient = snapClient
	k.storkClient = storkClient
	k.apiExtensionClient = apiExtensionClient
	k.dynamicInterface = dynamicInterface
	k.ocpClient = ocpClient
	k.ocpSecurityClient = ocpSecurityClient
	k.autopilotClient = autopilotClient

	k.setClients()
}

// SetBaseClient sets the kubernetes clientset
func (k *k8sOps) SetBaseClient(client kubernetes.Interface) {
	k.client = client
	k.setClients()
}

// SetSnapshotClient sets the snapshot clientset
func (k *k8sOps) SetSnapshotClient(snapClient rest.Interface) {
	k.snapClient = snapClient
	k.setClients()
}

// SetStorkClient sets the stork clientset
func (k *k8sOps) SetStorkClient(storkClient storkclientset.Interface) {
	k.storkClient = storkClient
	k.setClients()
}

// SetOpenstorageOperatorClient sets the openstorage operator clientset
func (k *k8sOps) SetOpenstorageOperatorClient(ostClient ostclientset.Interface) {
	k.ostClient = ostClient
	k.setClients()
}

// SetAPIExtensionsClient sets the api extensions clientset
func (k *k8sOps) SetAPIExtensionsClient(apiExtensionsClient apiextensionsclient.Interface) {
	k.apiExtensionClient = apiExtensionsClient
	k.setClients()
}

// SetDynamicClient sets the dynamic clientset
func (k *k8sOps) SetDynamicClient(dynamicClient dynamicclient.Interface) {
	k.dynamicInterface = dynamicClient
	k.setClients()
}

// SetOpenshiftAppsClient sets the openshift apps clientset
func (k *k8sOps) SetOpenshiftAppsClient(ocpAppsClient ocp_clientset.Interface) {
	k.ocpClient = ocpAppsClient
	k.setClients()
}

// SetOpenshiftSecurityClient sets the openshift security clientset
func (k *k8sOps) SetOpenshiftSecurityClient(ocpSecurityClient ocp_security_clientset.Interface) {
	k.ocpSecurityClient = ocpSecurityClient
	k.setClients()
}

// SetAutopilotClient sets the autopilot clientset
func (k *k8sOps) SetAutopilotClient(autopilotClient autopilotclientset.Interface) {
	k.autopilotClient = autopilotClient
	k.setClients()
}

// SetTalismanClient sets the talisman clientset
func (k *k8sOps) SetTalismanClient(talismanClient talismanclientset.Interface) {
	k.talismanClient = talismanClient
	k.setClients()
}

// SetPrometheusClient sets the prometheus clientset
func (k *k8sOps) SetPrometheusClient(prometheusClient prometheusclient.Interface) {
	k.prometheusClient = prometheusClient
	k.setClients()
}

func (k *k8sOps) setClients() {
	admissionregistration.SetInstance(admissionregistration.New(k.client.AdmissionregistrationV1beta1()))
	apiextensions.SetInstance(apiextensions.New(k.apiExtensionClient))
	apps.SetInstance(apps.New(k.client.AppsV1(), k.client.CoreV1()))
	autopilot.SetInstance(autopilot.New(k.autopilotClient))
	batch.SetInstance(batch.New(k.client.BatchV1()))
	core.SetInstance(core.New(k.client.CoreV1(), k.client.StorageV1()))
	discovery.SetInstance(discovery.New(k.client.Discovery()))
	dynamic.SetInstance(dynamic.New(k.dynamicInterface))
	externalstorage.SetInstance(externalstorage.New(k.snapClient))
	openshift.SetInstance(openshift.New(k.client, k.ocpClient, k.ocpSecurityClient))
	operator.SetInstance(operator.New(k.ostClient))
	prometheus.SetInstance(prometheus.New(k.prometheusClient))
	rbac.SetInstance(rbac.New(k.client.RbacV1()))
	storage.SetInstance(storage.New(k.client.StorageV1(), k.client.StorageV1beta1()))
	stork.SetInstance(stork.New(k.client, k.storkClient, k.snapClient))
	talisman.SetInstance(talisman.New(k.talismanClient))
}
