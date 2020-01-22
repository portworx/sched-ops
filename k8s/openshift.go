package k8s

import (
	"time"

	ocp_appsv1_api "github.com/openshift/api/apps/v1"
	ocp_securityv1_api "github.com/openshift/api/security/v1"
	"github.com/portworx/sched-ops/k8s/openshift"
	v1 "k8s.io/api/core/v1"
)

// Security Context Constraints APIs - BEGIN

func (k *k8sOps) ListSecurityContextConstraints() (result *ocp_securityv1_api.SecurityContextConstraintsList, err error) {
	return openshift.Instance().ListSecurityContextConstraints()
}

func (k *k8sOps) GetSecurityContextConstraints(name string) (result *ocp_securityv1_api.SecurityContextConstraints, err error) {
	return openshift.Instance().GetSecurityContextConstraints(name)
}

func (k *k8sOps) UpdateSecurityContextConstraints(securityContextConstraints *ocp_securityv1_api.SecurityContextConstraints) (result *ocp_securityv1_api.SecurityContextConstraints, err error) {
	return openshift.Instance().UpdateSecurityContextConstraints(securityContextConstraints)
}

// Security Context Constraints APIs - END

// DeploymentConfig APIs - BEGIN

func (k *k8sOps) ListDeploymentConfigs(namespace string) (*ocp_appsv1_api.DeploymentConfigList, error) {
	return openshift.Instance().ListDeploymentConfigs(namespace)
}

func (k *k8sOps) GetDeploymentConfig(name, namespace string) (*ocp_appsv1_api.DeploymentConfig, error) {
	return openshift.Instance().GetDeploymentConfig(name, namespace)
}

func (k *k8sOps) CreateDeploymentConfig(deployment *ocp_appsv1_api.DeploymentConfig) (*ocp_appsv1_api.DeploymentConfig, error) {
	return openshift.Instance().CreateDeploymentConfig(deployment)
}

func (k *k8sOps) DeleteDeploymentConfig(name, namespace string) error {
	return openshift.Instance().DeleteDeploymentConfig(name, namespace)
}

func (k *k8sOps) DescribeDeploymentConfig(depName, depNamespace string) (*ocp_appsv1_api.DeploymentConfigStatus, error) {
	return openshift.Instance().DescribeDeploymentConfig(depName, depNamespace)
}

func (k *k8sOps) UpdateDeploymentConfig(deployment *ocp_appsv1_api.DeploymentConfig) (*ocp_appsv1_api.DeploymentConfig, error) {
	return openshift.Instance().UpdateDeploymentConfig(deployment)
}

func (k *k8sOps) ValidateDeploymentConfig(deployment *ocp_appsv1_api.DeploymentConfig, timeout, retryInterval time.Duration) error {
	return openshift.Instance().ValidateDeploymentConfig(deployment, timeout, retryInterval)
}

func (k *k8sOps) ValidateTerminatedDeploymentConfig(deployment *ocp_appsv1_api.DeploymentConfig) error {
	return openshift.Instance().ValidateTerminatedDeploymentConfig(deployment)
}

func (k *k8sOps) GetDeploymentConfigPods(deployment *ocp_appsv1_api.DeploymentConfig) ([]v1.Pod, error) {
	return openshift.Instance().GetDeploymentConfigPods(deployment)
}

func (k *k8sOps) GetDeploymentConfigsUsingStorageClass(scName string) ([]ocp_appsv1_api.DeploymentConfig, error) {
	return openshift.Instance().GetDeploymentConfigsUsingStorageClass(scName)
}

// DeploymentConfig APIs - END
