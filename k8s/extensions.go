package k8s

import (
	"time"

	"github.com/portworx/sched-ops/k8s/apiextensions"
	apiextensionsv1beta1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
)

// CRD APIs - BEGIN

func (k *k8sOps) CreateCRD(resource apiextensions.CustomResource) error {
	return apiextensions.Instance().CreateCRD(resource)
}

func (k *k8sOps) RegisterCRD(crd *apiextensionsv1beta1.CustomResourceDefinition) error {
	return apiextensions.Instance().RegisterCRD(crd)
}

func (k *k8sOps) ValidateCRD(resource apiextensions.CustomResource, timeout, retryInterval time.Duration) error {
	return apiextensions.Instance().ValidateCRD(resource, timeout, retryInterval)
}

func (k *k8sOps) DeleteCRD(fullName string) error {
	return apiextensions.Instance().DeleteCRD(fullName)
}

// CRD APIs - END
