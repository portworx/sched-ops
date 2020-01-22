package k8s

import (
	"sync"
	"time"

	prometheusclient "github.com/coreos/prometheus-operator/pkg/client/versioned"
	snap_v1 "github.com/kubernetes-incubator/external-storage/snapshot/pkg/apis/crd/v1"
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
	hook "k8s.io/api/admissionregistration/v1beta1"
	appsv1 "k8s.io/api/apps/v1"
	batch_v1 "k8s.io/api/batch/v1"
	v1 "k8s.io/api/core/v1"
	rbac_v1 "k8s.io/api/rbac/v1"
	storagev1 "k8s.io/api/storage/v1"
	storagev1beta1 "k8s.io/api/storage/v1beta1"
	apiextensionsclient "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/version"
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
	discovery.Interface
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

var (
	instance Ops
	once     sync.Once
)

type k8sOps struct {
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
		instance = &k8sOps{}
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
	// set config for k8s clients
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

	return &k8sOps{}, nil
}

func (k *k8sOps) GetVersion() (*version.Info, error) {
	return discovery.Instance().GetVersion()
}

// Namespace APIs - BEGIN

func (k *k8sOps) ListNamespaces(labelSelector map[string]string) (*v1.NamespaceList, error) {
	return core.Instance().ListNamespaces(labelSelector)
}

func (k *k8sOps) GetNamespace(name string) (*v1.Namespace, error) {
	return core.Instance().GetNamespace(name)
}

func (k *k8sOps) CreateNamespace(name string, metadata map[string]string) (*v1.Namespace, error) {
	return core.Instance().CreateNamespace(name, metadata)
}

func (k *k8sOps) DeleteNamespace(name string) error {
	return core.Instance().DeleteNamespace(name)
}

// Namespace APIs - END
func (k *k8sOps) CreateNode(n *v1.Node) (*v1.Node, error) {
	return core.Instance().CreateNode(n)
}

func (k *k8sOps) UpdateNode(n *v1.Node) (*v1.Node, error) {
	return core.Instance().UpdateNode(n)
}

func (k *k8sOps) GetNodes() (*v1.NodeList, error) {
	return core.Instance().GetNodes()
}

func (k *k8sOps) GetNodeByName(name string) (*v1.Node, error) {
	return core.Instance().GetNodeByName(name)
}

func (k *k8sOps) IsNodeReady(name string) error {
	return core.Instance().IsNodeReady(name)
}

func (k *k8sOps) IsNodeMaster(node v1.Node) bool {
	return core.Instance().IsNodeMaster(node)
}

func (k *k8sOps) GetLabelsOnNode(name string) (map[string]string, error) {
	return core.Instance().GetLabelsOnNode(name)
}

// SearchNodeByAddresses searches the node based on the IP addresses, then it falls back to a
// search by hostname, and finally by the labels
func (k *k8sOps) SearchNodeByAddresses(addresses []string) (*v1.Node, error) {
	return core.Instance().SearchNodeByAddresses(addresses)
}

// FindMyNode finds LOCAL Node in Kubernetes cluster.
func (k *k8sOps) FindMyNode() (*v1.Node, error) {
	return core.Instance().FindMyNode()
}

func (k *k8sOps) AddLabelOnNode(name, key, value string) error {
	return core.Instance().AddLabelOnNode(name, key, value)
}

func (k *k8sOps) RemoveLabelOnNode(name, key string) error {
	return core.Instance().RemoveLabelOnNode(name, key)
}

func (k *k8sOps) WatchNode(node *v1.Node, watchNodeFn core.WatchFunc) error {
	return core.Instance().WatchNode(node, watchNodeFn)
}

func (k *k8sOps) CordonNode(nodeName string, timeout, retryInterval time.Duration) error {
	return core.Instance().CordonNode(nodeName, timeout, retryInterval)
}

func (k *k8sOps) UnCordonNode(nodeName string, timeout, retryInterval time.Duration) error {
	return core.Instance().UnCordonNode(nodeName, timeout, retryInterval)
}

func (k *k8sOps) DrainPodsFromNode(nodeName string, pods []v1.Pod, timeout time.Duration, retryInterval time.Duration) error {
	return core.Instance().DrainPodsFromNode(nodeName, pods, timeout, retryInterval)
}

func (k *k8sOps) WaitForPodDeletion(uid types.UID, namespace string, timeout time.Duration) error {
	return core.Instance().WaitForPodDeletion(uid, namespace, timeout)
}

func (k *k8sOps) RunCommandInPod(cmds []string, podName, containerName, namespace string) (string, error) {
	return core.Instance().RunCommandInPod(cmds, podName, containerName, namespace)
}

// Service APIs - BEGIN

func (k *k8sOps) CreateService(service *v1.Service) (*v1.Service, error) {
	return core.Instance().CreateService(service)
}

func (k *k8sOps) DeleteService(name, namespace string) error {
	return core.Instance().DeleteService(name, namespace)
}

func (k *k8sOps) GetService(svcName string, svcNS string) (*v1.Service, error) {
	return core.Instance().GetService(svcName, svcNS)
}

func (k *k8sOps) DescribeService(svcName string, svcNamespace string) (*v1.ServiceStatus, error) {
	return core.Instance().DescribeService(svcName, svcNamespace)
}

func (k *k8sOps) ValidateDeletedService(svcName string, svcNS string) error {
	return core.Instance().ValidateDeletedService(svcName, svcNS)
}

func (k *k8sOps) PatchService(name, namespace string, jsonPatch []byte) (*v1.Service, error) {
	return core.Instance().PatchService(name, namespace, jsonPatch)
}

// Service APIs - END

// Deployment APIs - BEGIN

func (k *k8sOps) ListDeployments(namespace string, options meta_v1.ListOptions) (*appsv1.DeploymentList, error) {
	return apps.Instance().ListDeployments(namespace, options)
}

func (k *k8sOps) GetDeployment(name, namespace string) (*appsv1.Deployment, error) {
	return apps.Instance().GetDeployment(name, namespace)
}

func (k *k8sOps) CreateDeployment(deployment *appsv1.Deployment) (*appsv1.Deployment, error) {
	return apps.Instance().CreateDeployment(deployment)
}

func (k *k8sOps) DeleteDeployment(name, namespace string) error {
	return apps.Instance().DeleteDeployment(name, namespace)
}

func (k *k8sOps) DescribeDeployment(depName, depNamespace string) (*appsv1.DeploymentStatus, error) {
	return apps.Instance().DescribeDeployment(depName, depNamespace)
}

func (k *k8sOps) UpdateDeployment(deployment *appsv1.Deployment) (*appsv1.Deployment, error) {
	return apps.Instance().UpdateDeployment(deployment)
}

func (k *k8sOps) ValidateDeployment(deployment *appsv1.Deployment, timeout, retryInterval time.Duration) error {
	return apps.Instance().ValidateDeployment(deployment, timeout, retryInterval)
}

func (k *k8sOps) ValidateTerminatedDeployment(deployment *appsv1.Deployment, timeout, timeBeforeRetry time.Duration) error {
	return apps.Instance().ValidateTerminatedDeployment(deployment, timeout, timeBeforeRetry)
}

func (k *k8sOps) GetDeploymentPods(deployment *appsv1.Deployment) ([]v1.Pod, error) {
	return apps.Instance().GetDeploymentPods(deployment)
}

func (k *k8sOps) GetDeploymentsUsingStorageClass(scName string) ([]appsv1.Deployment, error) {
	return apps.Instance().GetDeploymentsUsingStorageClass(scName)
}

// Deployment APIs - END

// DaemonSet APIs - BEGIN

func (k *k8sOps) CreateDaemonSet(ds *appsv1.DaemonSet) (*appsv1.DaemonSet, error) {
	return apps.Instance().CreateDaemonSet(ds)
}

func (k *k8sOps) ListDaemonSets(namespace string, listOpts meta_v1.ListOptions) ([]appsv1.DaemonSet, error) {
	return apps.Instance().ListDaemonSets(namespace, listOpts)
}

func (k *k8sOps) GetDaemonSet(name, namespace string) (*appsv1.DaemonSet, error) {
	return apps.Instance().GetDaemonSet(name, namespace)
}

func (k *k8sOps) GetDaemonSetPods(ds *appsv1.DaemonSet) ([]v1.Pod, error) {
	return apps.Instance().GetDaemonSetPods(ds)
}

func (k *k8sOps) ValidateDaemonSet(name, namespace string, timeout time.Duration) error {
	return apps.Instance().ValidateDaemonSet(name, namespace, timeout)
}

func (k *k8sOps) UpdateDaemonSet(ds *appsv1.DaemonSet) (*appsv1.DaemonSet, error) {
	return apps.Instance().UpdateDaemonSet(ds)
}

func (k *k8sOps) DeleteDaemonSet(name, namespace string) error {
	return apps.Instance().DeleteDaemonSet(name, namespace)
}

// DaemonSet APIs - END

// Job APIs - BEGIN
func (k *k8sOps) CreateJob(job *batch_v1.Job) (*batch_v1.Job, error) {
	return batch.Instance().CreateJob(job)
}

func (k *k8sOps) GetJob(name, namespace string) (*batch_v1.Job, error) {
	return batch.Instance().GetJob(name, namespace)
}

func (k *k8sOps) DeleteJob(name, namespace string) error {
	return batch.Instance().DeleteJob(name, namespace)
}

func (k *k8sOps) ValidateJob(name, namespace string, timeout time.Duration) error {
	return batch.Instance().ValidateJob(name, namespace, timeout)
}

// Job APIs - END

// StatefulSet APIs - BEGIN

func (k *k8sOps) ListStatefulSets(namespace string) (*appsv1.StatefulSetList, error) {
	return apps.Instance().ListStatefulSets(namespace)
}

func (k *k8sOps) GetStatefulSet(name, namespace string) (*appsv1.StatefulSet, error) {
	return apps.Instance().GetStatefulSet(name, namespace)
}

func (k *k8sOps) CreateStatefulSet(statefulset *appsv1.StatefulSet) (*appsv1.StatefulSet, error) {
	return apps.Instance().CreateStatefulSet(statefulset)
}

func (k *k8sOps) DeleteStatefulSet(name, namespace string) error {
	return apps.Instance().DeleteStatefulSet(name, namespace)
}

func (k *k8sOps) DescribeStatefulSet(ssetName string, ssetNamespace string) (*appsv1.StatefulSetStatus, error) {
	return apps.Instance().DescribeStatefulSet(ssetName, ssetNamespace)
}

func (k *k8sOps) UpdateStatefulSet(statefulset *appsv1.StatefulSet) (*appsv1.StatefulSet, error) {
	return apps.Instance().UpdateStatefulSet(statefulset)
}

func (k *k8sOps) ValidateStatefulSet(statefulset *appsv1.StatefulSet, timeout time.Duration) error {
	return apps.Instance().ValidateStatefulSet(statefulset, timeout)
}

func (k *k8sOps) GetStatefulSetPods(statefulset *appsv1.StatefulSet) ([]v1.Pod, error) {
	return apps.Instance().GetStatefulSetPods(statefulset)
}

func (k *k8sOps) ValidateTerminatedStatefulSet(statefulset *appsv1.StatefulSet, timeout, timeBeforeRetry time.Duration) error {
	return apps.Instance().ValidateTerminatedStatefulSet(statefulset, timeout, timeBeforeRetry)
}

func (k *k8sOps) GetStatefulSetsUsingStorageClass(scName string) ([]appsv1.StatefulSet, error) {
	return apps.Instance().GetStatefulSetsUsingStorageClass(scName)
}

func (k *k8sOps) GetPVCsForStatefulSet(ss *appsv1.StatefulSet) (*v1.PersistentVolumeClaimList, error) {
	return apps.Instance().GetPVCsForStatefulSet(ss)
}

func (k *k8sOps) ValidatePVCsForStatefulSet(ss *appsv1.StatefulSet, timeout, retryTimeout time.Duration) error {
	return apps.Instance().ValidatePVCsForStatefulSet(ss, timeout, retryTimeout)
}

// StatefulSet APIs - END

// RBAC APIs - BEGIN

func (k *k8sOps) CreateRole(role *rbac_v1.Role) (*rbac_v1.Role, error) {
	return rbac.Instance().CreateRole(role)
}

func (k *k8sOps) UpdateRole(role *rbac_v1.Role) (*rbac_v1.Role, error) {
	return rbac.Instance().UpdateRole(role)
}

func (k *k8sOps) GetRole(name, namespace string) (*rbac_v1.Role, error) {
	return rbac.Instance().GetRole(name, namespace)
}

func (k *k8sOps) CreateClusterRole(role *rbac_v1.ClusterRole) (*rbac_v1.ClusterRole, error) {
	return rbac.Instance().CreateClusterRole(role)
}

func (k *k8sOps) GetClusterRole(name string) (*rbac_v1.ClusterRole, error) {
	return rbac.Instance().GetClusterRole(name)
}

func (k *k8sOps) UpdateClusterRole(role *rbac_v1.ClusterRole) (*rbac_v1.ClusterRole, error) {
	return rbac.Instance().UpdateClusterRole(role)
}

func (k *k8sOps) CreateRoleBinding(binding *rbac_v1.RoleBinding) (*rbac_v1.RoleBinding, error) {
	return rbac.Instance().CreateRoleBinding(binding)
}

func (k *k8sOps) UpdateRoleBinding(binding *rbac_v1.RoleBinding) (*rbac_v1.RoleBinding, error) {
	return rbac.Instance().UpdateRoleBinding(binding)
}

func (k *k8sOps) GetRoleBinding(name, namespace string) (*rbac_v1.RoleBinding, error) {
	return rbac.Instance().GetRoleBinding(name, namespace)
}

func (k *k8sOps) CreateClusterRoleBinding(binding *rbac_v1.ClusterRoleBinding) (*rbac_v1.ClusterRoleBinding, error) {
	return rbac.Instance().CreateClusterRoleBinding(binding)
}

func (k *k8sOps) UpdateClusterRoleBinding(binding *rbac_v1.ClusterRoleBinding) (*rbac_v1.ClusterRoleBinding, error) {
	return rbac.Instance().UpdateClusterRoleBinding(binding)
}

func (k *k8sOps) GetClusterRoleBinding(name string) (*rbac_v1.ClusterRoleBinding, error) {
	return rbac.Instance().GetClusterRoleBinding(name)
}

func (k *k8sOps) ListClusterRoleBindings() (*rbac_v1.ClusterRoleBindingList, error) {
	return rbac.Instance().ListClusterRoleBindings()
}

func (k *k8sOps) CreateServiceAccount(account *v1.ServiceAccount) (*v1.ServiceAccount, error) {
	return core.Instance().CreateServiceAccount(account)
}

func (k *k8sOps) GetServiceAccount(name, namespace string) (*v1.ServiceAccount, error) {
	return core.Instance().GetServiceAccount(name, namespace)
}

func (k *k8sOps) DeleteRole(name, namespace string) error {
	return rbac.Instance().DeleteRole(name, namespace)
}

func (k *k8sOps) DeleteClusterRole(roleName string) error {
	return rbac.Instance().DeleteClusterRole(roleName)
}

func (k *k8sOps) DeleteRoleBinding(name, namespace string) error {
	return rbac.Instance().DeleteRoleBinding(name, namespace)
}

func (k *k8sOps) DeleteClusterRoleBinding(bindingName string) error {
	return rbac.Instance().DeleteClusterRoleBinding(bindingName)
}

func (k *k8sOps) DeleteServiceAccount(accountName, namespace string) error {
	return core.Instance().DeleteServiceAccount(accountName, namespace)
}

// RBAC APIs - END

// Pod APIs - BEGIN

func (k *k8sOps) DeletePods(pods []v1.Pod, force bool) error {
	return core.Instance().DeletePods(pods, force)
}

func (k *k8sOps) DeletePod(name string, ns string, force bool) error {
	return core.Instance().DeletePod(name, ns, force)
}

func (k *k8sOps) CreatePod(pod *v1.Pod) (*v1.Pod, error) {
	return core.Instance().CreatePod(pod)
}

func (k *k8sOps) UpdatePod(pod *v1.Pod) (*v1.Pod, error) {
	return core.Instance().UpdatePod(pod)
}

func (k *k8sOps) GetPods(namespace string, labelSelector map[string]string) (*v1.PodList, error) {
	return core.Instance().GetPods(namespace, labelSelector)
}

func (k *k8sOps) GetPodsByNode(nodeName, namespace string) (*v1.PodList, error) {
	return core.Instance().GetPodsByNode(nodeName, namespace)
}

func (k *k8sOps) GetPodsByOwner(ownerUID types.UID, namespace string) ([]v1.Pod, error) {
	return core.Instance().GetPodsByOwner(ownerUID, namespace)
}

func (k *k8sOps) GetPodsUsingPV(pvName string) ([]v1.Pod, error) {
	return core.Instance().GetPodsUsingPV(pvName)
}

func (k *k8sOps) GetPodsUsingPVByNodeName(pvName, nodeName string) ([]v1.Pod, error) {
	return core.Instance().GetPodsUsingPVByNodeName(pvName, nodeName)
}

func (k *k8sOps) GetPodsUsingPVC(pvcName, pvcNamespace string) ([]v1.Pod, error) {
	return core.Instance().GetPodsUsingPVC(pvcName, pvcNamespace)
}

func (k *k8sOps) GetPodsUsingPVCByNodeName(pvcName, pvcNamespace, nodeName string) ([]v1.Pod, error) {
	return core.Instance().GetPodsUsingPVCByNodeName(pvcName, pvcNamespace, nodeName)
}

func (k *k8sOps) GetPodsUsingVolumePlugin(plugin string) ([]v1.Pod, error) {
	return core.Instance().GetPodsUsingVolumePlugin(plugin)
}

func (k *k8sOps) GetPodsUsingVolumePluginByNodeName(nodeName, plugin string) ([]v1.Pod, error) {
	return core.Instance().GetPodsUsingVolumePluginByNodeName(nodeName, plugin)
}

func (k *k8sOps) GetPodByName(podName string, namespace string) (*v1.Pod, error) {
	return core.Instance().GetPodByName(podName, namespace)
}

func (k *k8sOps) GetPodByUID(uid types.UID, namespace string) (*v1.Pod, error) {
	return core.Instance().GetPodByUID(uid, namespace)
}

func (k *k8sOps) IsPodRunning(pod v1.Pod) bool {
	return core.Instance().IsPodRunning(pod)
}

func (k *k8sOps) IsPodReady(pod v1.Pod) bool {
	return core.Instance().IsPodReady(pod)
}

func (k *k8sOps) IsPodBeingManaged(pod v1.Pod) bool {
	return core.Instance().IsPodBeingManaged(pod)
}

func (k *k8sOps) ValidatePod(pod *v1.Pod, timeout, retryInterval time.Duration) error {
	return core.Instance().ValidatePod(pod, timeout, retryInterval)
}

func (k *k8sOps) WatchPods(namespace string, fn core.WatchFunc, listOptions meta_v1.ListOptions) error {
	return core.Instance().WatchPods(namespace, fn, listOptions)
}

// Pod APIs - END

// StorageClass APIs - BEGIN

func (k *k8sOps) GetStorageClasses(labelSelector map[string]string) (*storagev1.StorageClassList, error) {
	return storage.Instance().GetStorageClasses(labelSelector)
}

func (k *k8sOps) GetStorageClass(name string) (*storagev1.StorageClass, error) {
	return storage.Instance().GetStorageClass(name)
}

func (k *k8sOps) CreateStorageClass(sc *storagev1.StorageClass) (*storagev1.StorageClass, error) {
	return storage.Instance().CreateStorageClass(sc)
}

func (k *k8sOps) DeleteStorageClass(name string) error {
	return storage.Instance().DeleteStorageClass(name)
}

func (k *k8sOps) GetStorageClassParams(sc *storagev1.StorageClass) (map[string]string, error) {
	return storage.Instance().GetStorageClassParams(sc)
}

func (k *k8sOps) ValidateStorageClass(name string) (*storagev1.StorageClass, error) {
	return storage.Instance().ValidateStorageClass(name)
}

// StorageClass APIs - END

// PVC APIs - BEGIN

func (k *k8sOps) CreatePersistentVolumeClaim(pvc *v1.PersistentVolumeClaim) (*v1.PersistentVolumeClaim, error) {
	return core.Instance().CreatePersistentVolumeClaim(pvc)
}

func (k *k8sOps) UpdatePersistentVolumeClaim(pvc *v1.PersistentVolumeClaim) (*v1.PersistentVolumeClaim, error) {
	return core.Instance().UpdatePersistentVolumeClaim(pvc)
}

func (k *k8sOps) DeletePersistentVolumeClaim(name, namespace string) error {
	return core.Instance().DeletePersistentVolumeClaim(name, namespace)
}

func (k *k8sOps) ValidatePersistentVolumeClaim(pvc *v1.PersistentVolumeClaim, timeout, retryInterval time.Duration) error {
	return core.Instance().ValidatePersistentVolumeClaim(pvc, timeout, retryInterval)
}

func (k *k8sOps) ValidatePersistentVolumeClaimSize(pvc *v1.PersistentVolumeClaim, expectedPVCSize int64, timeout, retryInterval time.Duration) error {
	return core.Instance().ValidatePersistentVolumeClaimSize(pvc, expectedPVCSize, timeout, retryInterval)
}

func (k *k8sOps) CreatePersistentVolume(pv *v1.PersistentVolume) (*v1.PersistentVolume, error) {
	return core.Instance().CreatePersistentVolume(pv)
}

func (k *k8sOps) GetPersistentVolumeClaim(pvcName string, namespace string) (*v1.PersistentVolumeClaim, error) {
	return core.Instance().GetPersistentVolumeClaim(pvcName, namespace)
}

func (k *k8sOps) GetPersistentVolumeClaims(namespace string, labelSelector map[string]string) (*v1.PersistentVolumeClaimList, error) {
	return core.Instance().GetPersistentVolumeClaims(namespace, labelSelector)
}

func (k *k8sOps) GetPersistentVolume(pvName string) (*v1.PersistentVolume, error) {
	return core.Instance().GetPersistentVolume(pvName)
}

func (k *k8sOps) DeletePersistentVolume(pvName string) error {
	return core.Instance().DeletePersistentVolume(pvName)
}

func (k *k8sOps) GetPersistentVolumes() (*v1.PersistentVolumeList, error) {
	return core.Instance().GetPersistentVolumes()
}

func (k *k8sOps) GetVolumeForPersistentVolumeClaim(pvc *v1.PersistentVolumeClaim) (string, error) {
	return core.Instance().GetVolumeForPersistentVolumeClaim(pvc)
}

func (k *k8sOps) GetPersistentVolumeClaimStatus(pvc *v1.PersistentVolumeClaim) (*v1.PersistentVolumeClaimStatus, error) {
	return core.Instance().GetPersistentVolumeClaimStatus(pvc)
}

func (k *k8sOps) GetPersistentVolumeClaimParams(pvc *v1.PersistentVolumeClaim) (map[string]string, error) {
	return core.Instance().GetPersistentVolumeClaimParams(pvc)
}

func (k *k8sOps) GetPVCsUsingStorageClass(scName string) ([]v1.PersistentVolumeClaim, error) {
	return core.Instance().GetPVCsUsingStorageClass(scName)
}

func (k *k8sOps) GetStorageProvisionerForPVC(pvc *v1.PersistentVolumeClaim) (string, error) {
	return core.Instance().GetStorageProvisionerForPVC(pvc)
}

// PVCs APIs - END

// Snapshot APIs - BEGIN

func (k *k8sOps) CreateSnapshot(snap *snap_v1.VolumeSnapshot) (*snap_v1.VolumeSnapshot, error) {
	return externalstorage.Instance().CreateSnapshot(snap)
}

func (k *k8sOps) UpdateSnapshot(snap *snap_v1.VolumeSnapshot) (*snap_v1.VolumeSnapshot, error) {
	return externalstorage.Instance().UpdateSnapshot(snap)
}

func (k *k8sOps) DeleteSnapshot(name string, namespace string) error {
	return externalstorage.Instance().DeleteSnapshot(name, namespace)
}

func (k *k8sOps) ValidateSnapshot(name string, namespace string, retry bool, timeout, retryInterval time.Duration) error {
	return externalstorage.Instance().ValidateSnapshot(name, namespace, retry, timeout, retryInterval)
}

func (k *k8sOps) ValidateSnapshotData(name string, retry bool, timeout, retryInterval time.Duration) error {
	return externalstorage.Instance().ValidateSnapshotData(name, retry, timeout, retryInterval)
}

func (k *k8sOps) GetVolumeForSnapshot(name string, namespace string) (string, error) {
	return externalstorage.Instance().GetVolumeForSnapshot(name, namespace)
}

func (k *k8sOps) GetSnapshot(name string, namespace string) (*snap_v1.VolumeSnapshot, error) {
	return externalstorage.Instance().GetSnapshot(name, namespace)
}

func (k *k8sOps) ListSnapshots(namespace string) (*snap_v1.VolumeSnapshotList, error) {
	return externalstorage.Instance().ListSnapshots(namespace)
}

func (k *k8sOps) GetSnapshotStatus(name string, namespace string) (*snap_v1.VolumeSnapshotStatus, error) {
	return externalstorage.Instance().GetSnapshotStatus(name, namespace)
}

func (k *k8sOps) GetSnapshotData(name string) (*snap_v1.VolumeSnapshotData, error) {
	return externalstorage.Instance().GetSnapshotData(name)
}

func (k *k8sOps) CreateSnapshotData(snapData *snap_v1.VolumeSnapshotData) (*snap_v1.VolumeSnapshotData, error) {
	return externalstorage.Instance().CreateSnapshotData(snapData)
}

func (k *k8sOps) DeleteSnapshotData(name string) error {
	return externalstorage.Instance().DeleteSnapshotData(name)
}

// Snapshot APIs - END

// Secret APIs - BEGIN

func (k *k8sOps) GetSecret(name string, namespace string) (*v1.Secret, error) {
	return core.Instance().GetSecret(name, namespace)
}

func (k *k8sOps) CreateSecret(secret *v1.Secret) (*v1.Secret, error) {
	return core.Instance().CreateSecret(secret)
}

func (k *k8sOps) UpdateSecret(secret *v1.Secret) (*v1.Secret, error) {
	return core.Instance().UpdateSecret(secret)
}

func (k *k8sOps) UpdateSecretData(name string, ns string, data map[string][]byte) (*v1.Secret, error) {
	return core.Instance().UpdateSecretData(name, ns, data)
}

func (k *k8sOps) DeleteSecret(name, namespace string) error {
	return core.Instance().DeleteSecret(name, namespace)
}

// Secret APIs - END

// ConfigMap APIs - BEGIN

func (k *k8sOps) GetConfigMap(name string, namespace string) (*v1.ConfigMap, error) {
	return core.Instance().GetConfigMap(name, namespace)
}

func (k *k8sOps) CreateConfigMap(configMap *v1.ConfigMap) (*v1.ConfigMap, error) {
	return core.Instance().CreateConfigMap(configMap)
}

func (k *k8sOps) DeleteConfigMap(name, namespace string) error {
	return core.Instance().DeleteConfigMap(name, namespace)
}

func (k *k8sOps) UpdateConfigMap(configMap *v1.ConfigMap) (*v1.ConfigMap, error) {
	return core.Instance().UpdateConfigMap(configMap)
}

func (k *k8sOps) WatchConfigMap(configMap *v1.ConfigMap, fn core.WatchFunc) error {
	return core.Instance().WatchConfigMap(configMap, fn)
}

// ConfigMap APIs - END

// Event APIs - BEGIN
// CreateEvent puts an event into k8s etcd
func (k *k8sOps) CreateEvent(event *v1.Event) (*v1.Event, error) {
	return core.Instance().CreateEvent(event)
}

// ListEvents retrieves all events registered with kubernetes
func (k *k8sOps) ListEvents(namespace string, opts meta_v1.ListOptions) (*v1.EventList, error) {
	return core.Instance().ListEvents(namespace, opts)
}

// Event APIs - END

// Object APIs - BEGIN

// GetObject returns the latest object given a generic Object
func (k *k8sOps) GetObject(object runtime.Object) (runtime.Object, error) {
	return dynamic.Instance().GetObject(object)
}

// UpdateObject updates a generic Object
func (k *k8sOps) UpdateObject(object runtime.Object) (runtime.Object, error) {
	return dynamic.Instance().UpdateObject(object)
}

// Object APIs - END

// VolumeAttachment APIs - START

func (k *k8sOps) ListVolumeAttachments() (*storagev1beta1.VolumeAttachmentList, error) {
	return storage.Instance().ListVolumeAttachments()
}

func (k *k8sOps) DeleteVolumeAttachment(name string) error {
	return storage.Instance().DeleteVolumeAttachment(name)
}

func (k *k8sOps) CreateVolumeAttachment(volumeAttachment *storagev1beta1.VolumeAttachment) (*storagev1beta1.VolumeAttachment, error) {
	return storage.Instance().CreateVolumeAttachment(volumeAttachment)
}

func (k *k8sOps) UpdateVolumeAttachment(volumeAttachment *storagev1beta1.VolumeAttachment) (*storagev1beta1.VolumeAttachment, error) {
	return storage.Instance().UpdateVolumeAttachment(volumeAttachment)
}

func (k *k8sOps) UpdateVolumeAttachmentStatus(volumeAttachment *storagev1beta1.VolumeAttachment) (*storagev1beta1.VolumeAttachment, error) {
	return storage.Instance().UpdateVolumeAttachmentStatus(volumeAttachment)
}

// VolumeAttachment APIs - END

// MutatingWebhookConfig APIS - START

// GetMutatingWebhookConfiguration returns a given MutatingWebhookConfiguration
func (k *k8sOps) GetMutatingWebhookConfiguration(name string) (*hook.MutatingWebhookConfiguration, error) {
	return admissionregistration.Instance().GetMutatingWebhookConfiguration(name)
}

// CreateMutatingWebhookConfiguration creates given MutatingWebhookConfiguration
func (k *k8sOps) CreateMutatingWebhookConfiguration(cfg *hook.MutatingWebhookConfiguration) (*hook.MutatingWebhookConfiguration, error) {
	return admissionregistration.Instance().CreateMutatingWebhookConfiguration(cfg)
}

// UpdateMutatingWebhookConfiguration updates given MutatingWebhookConfiguration
func (k *k8sOps) UpdateMutatingWebhookConfiguration(cfg *hook.MutatingWebhookConfiguration) (*hook.MutatingWebhookConfiguration, error) {
	return admissionregistration.Instance().UpdateMutatingWebhookConfiguration(cfg)
}

// DeleteMutatingWebhookConfiguration deletes given MutatingWebhookConfiguration
func (k *k8sOps) DeleteMutatingWebhookConfiguration(name string) error {
	return admissionregistration.Instance().DeleteMutatingWebhookConfiguration(name)
}

// MutatingWebhookConfig APIS - END
