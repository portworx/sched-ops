module github.com/portworx/sched-ops

go 1.12

require (
	github.com/kubernetes-csi/external-snapshotter/client/v4 v4.0.0
	github.com/kubernetes-incubator/external-storage v0.20.4-openstorage-rc2
	github.com/libopenstorage/autopilot-api v1.3.0
	github.com/libopenstorage/operator v0.0.0-20210303221358-0bb211a9908c
	github.com/libopenstorage/stork v1.4.1-0.20210903185636-5a1f8a4142bf
	github.com/openshift/api v0.0.0-20210105115604-44119421ec6b
	github.com/openshift/client-go v0.0.0-20210112165513-ebc401615f47
	// TODO: Vendor from pb-1874 branch. Need to change it to master.
	github.com/portworx/kdmp v0.4.1-0.20210929095959-dbee7f98b06f
	github.com/portworx/talisman v0.0.0-20210302012732-8af4564777f7
	github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring v0.46.0
	github.com/prometheus-operator/prometheus-operator/pkg/client v0.46.0
	github.com/sirupsen/logrus v1.8.1
	github.com/stretchr/testify v1.7.0
	k8s.io/api v0.20.4
	k8s.io/apiextensions-apiserver v0.20.4
	k8s.io/apimachinery v0.20.4
	k8s.io/client-go v12.0.0+incompatible
)

replace (
	github.com/kubernetes-incubator/external-storage => github.com/libopenstorage/external-storage v0.20.4-openstorage-rc2
	github.com/libopenstorage/autopilot-api => github.com/libopenstorage/autopilot-api v0.6.1-0.20210301232050-ca2633c6e114

	k8s.io/api => k8s.io/api v0.20.4
	k8s.io/apiextensions-apiserver => k8s.io/apiextensions-apiserver v0.20.4
	k8s.io/apimachinery => k8s.io/apimachinery v0.20.4
	k8s.io/apiserver => k8s.io/apiserver v0.20.4
	k8s.io/cli-runtime => k8s.io/cli-runtime v0.20.4
	k8s.io/client-go => k8s.io/client-go v0.20.4
	k8s.io/cloud-provider => k8s.io/cloud-provider v0.20.4
	k8s.io/cluster-bootstrap => k8s.io/cluster-bootstrap v0.20.4
	k8s.io/code-generator => k8s.io/code-generator v0.20.4
	k8s.io/component-base => k8s.io/component-base v0.20.4
	k8s.io/component-helpers => k8s.io/component-helpers v0.20.4
	k8s.io/controller-manager => k8s.io/controller-manager v0.20.4
	k8s.io/cri-api => k8s.io/cri-api v0.20.4
	k8s.io/csi-translation-lib => k8s.io/csi-translation-lib v0.20.4
	k8s.io/kube-aggregator => k8s.io/kube-aggregator v0.20.4
	k8s.io/kube-controller-manager => k8s.io/kube-controller-manager v0.20.4
	k8s.io/kube-proxy => k8s.io/kube-proxy v0.20.4
	k8s.io/kube-scheduler => k8s.io/kube-scheduler v0.20.4
	k8s.io/kubectl => k8s.io/kubectl v0.20.4
	k8s.io/kubelet => k8s.io/kubelet v0.20.4
	k8s.io/kubernetes => k8s.io/kubernetes v1.20.4
	k8s.io/legacy-cloud-providers => k8s.io/legacy-cloud-providers v0.20.4
	k8s.io/metrics => k8s.io/metrics v0.20.4
	k8s.io/mount-utils => k8s.io/mount-utils v0.20.4
	k8s.io/sample-apiserver => k8s.io/sample-apiserver v0.20.4
)
