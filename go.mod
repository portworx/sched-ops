module github.com/portworx/sched-ops

go 1.24.0

toolchain go1.24.4

require (
	github.com/libopenstorage/autopilot-api v1.3.0
	github.com/libopenstorage/openstorage v9.4.47+incompatible
	github.com/libopenstorage/operator v0.0.0-20230801044606-e27dec4914d4
	github.com/libopenstorage/stork v1.4.1-0.20230610103146-72cf75320066
	github.com/openshift/api v0.0.0-20230503133300-8bbcb7ca7183
	github.com/openshift/client-go v0.0.0-20210112165513-ebc401615f47
	// TODO: Vendor from pb-1874 branch. Need to change it to master.
	github.com/portworx/kdmp v0.4.1-0.20230830195819-704706dda5c7
	github.com/portworx/talisman v0.0.0-20210302012732-8af4564777f7
	github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring v0.68.0
	github.com/prometheus-operator/prometheus-operator/pkg/client v0.46.0
	github.com/sirupsen/logrus v1.9.3
	github.com/stretchr/testify v1.10.0
	k8s.io/api v0.33.5
	k8s.io/apiextensions-apiserver v0.33.5
	k8s.io/apimachinery v0.33.5
	k8s.io/client-go v12.0.0+incompatible
)

require (
	github.com/coreos/go-oidc v2.3.0+incompatible // indirect
	github.com/davecgh/go-spew v1.1.2-0.20180830191138-d8f796af33cc // indirect
	github.com/dgrijalva/jwt-go v3.2.1-0.20180719211823-0b96aaa70776+incompatible // indirect
	github.com/emicklei/go-restful/v3 v3.12.2 // indirect
	github.com/go-logr/logr v1.4.2 // indirect
	github.com/go-openapi/jsonpointer v0.21.0 // indirect
	github.com/go-openapi/jsonreference v0.20.2 // indirect
	github.com/go-openapi/swag v0.23.0 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/golang/groupcache v0.0.0-20210331224755-41bb18bfe9da // indirect
	github.com/golang/protobuf v1.5.4 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/grpc-ecosystem/go-grpc-middleware v1.3.0 // indirect
	github.com/grpc-ecosystem/grpc-gateway v1.16.0 // indirect
	github.com/josharian/intern v1.0.0 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/mailru/easyjson v0.7.7 // indirect
	github.com/moby/spdystream v0.5.0 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/mohae/deepcopy v0.0.0-20170929034955-c48cc78d4826 // indirect
	github.com/munnerz/goautoneg v0.0.0-20191010083416-a7dc8b61c822 // indirect
	github.com/operator-framework/api v0.17.1
	github.com/pborman/uuid v1.2.1 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/pmezard/go-difflib v1.0.1-0.20181226105442-5d4384ee4fb2 // indirect
	github.com/pquerna/cachecontrol v0.1.0 // indirect
	github.com/spf13/pflag v1.0.6 // indirect
	golang.org/x/crypto v0.36.0 // indirect
	golang.org/x/net v0.38.0 // indirect
	golang.org/x/oauth2 v0.28.0 // indirect
	golang.org/x/sys v0.31.0 // indirect
	golang.org/x/term v0.30.0 // indirect
	golang.org/x/text v0.23.0 // indirect
	golang.org/x/time v0.9.0 // indirect
	google.golang.org/genproto v0.0.0-20240227224415-6ceb2ff114de // indirect
	google.golang.org/grpc v1.68.1
	google.golang.org/protobuf v1.36.5 // indirect
	gopkg.in/inf.v0 v0.9.1 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

require (
	github.com/blang/semver/v4 v4.0.0 // indirect
	github.com/evanphx/json-patch/v5 v5.9.11 // indirect
)

require github.com/golang/glog v1.2.2 // indirect

require (
	github.com/kubernetes-csi/external-snapshotter/client/v6 v6.2.0
	github.com/kubernetes-incubator/external-storage v0.20.4-openstorage-rc7
	github.com/tektoncd/pipeline v0.56.0
	github.com/undefinedlabs/go-mpatch v1.0.7
	google.golang.org/api v0.169.0
	sigs.k8s.io/cluster-api v1.10.1
	sigs.k8s.io/controller-runtime v0.20.4
)

require (
	cel.dev/expr v0.19.1 // indirect
	cloud.google.com/go/compute/metadata v0.5.0 // indirect
	contrib.go.opencensus.io/exporter/ocagent v0.7.1-0.20200907061046-05415f1de66d // indirect
	contrib.go.opencensus.io/exporter/prometheus v0.4.0 // indirect
	github.com/antlr4-go/antlr/v4 v4.13.0 // indirect
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/blendle/zapdriver v1.3.1 // indirect
	github.com/census-instrumentation/opencensus-proto v0.4.1 // indirect
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/felixge/httpsnoop v1.0.4 // indirect
	github.com/fxamacker/cbor/v2 v2.7.0 // indirect
	github.com/go-kit/log v0.2.1 // indirect
	github.com/go-logfmt/logfmt v0.6.0 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/google/cel-go v0.23.2 // indirect
	github.com/google/gnostic-models v0.6.9 // indirect
	github.com/google/go-cmp v0.7.0 // indirect
	github.com/google/go-containerregistry v0.17.0 // indirect
	github.com/google/s2a-go v0.1.7 // indirect
	github.com/googleapis/enterprise-certificate-proxy v0.3.2 // indirect
	github.com/googleapis/gax-go/v2 v2.12.2 // indirect
	github.com/gorilla/websocket v1.5.4-0.20250319132907-e064f32e3674 // indirect
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.24.0 // indirect
	github.com/hashicorp/errwrap v1.1.0 // indirect
	github.com/hashicorp/go-multierror v1.1.1 // indirect
	github.com/k8snetworkplumbingwg/network-attachment-definition-client v0.0.0-20191119172530-79f836b90111 // indirect
	github.com/kubernetes-csi/external-snapshotter/client/v4 v4.2.0 // indirect
	github.com/mxk/go-flowrate v0.0.0-20140419014527-cca7078d478f // indirect
	github.com/opencontainers/go-digest v1.0.0 // indirect
	github.com/openshift/custom-resource-status v1.1.2 // indirect
	github.com/prometheus/client_golang v1.22.0 // indirect
	github.com/prometheus/client_model v0.6.1 // indirect
	github.com/prometheus/common v0.62.0 // indirect
	github.com/prometheus/procfs v0.15.1 // indirect
	github.com/prometheus/statsd_exporter v0.21.0 // indirect
	github.com/stoewer/go-strcase v1.3.0 // indirect
	github.com/x448/float16 v0.8.4 // indirect
	go.opencensus.io v0.24.0 // indirect
	go.opentelemetry.io/auto/sdk v1.1.0 // indirect
	go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp v0.58.0 // indirect
	go.opentelemetry.io/otel v1.33.0 // indirect
	go.opentelemetry.io/otel/metric v1.33.0 // indirect
	go.opentelemetry.io/otel/trace v1.33.0 // indirect
	go.uber.org/atomic v1.10.0 // indirect
	go.uber.org/mock v0.5.1 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	go.uber.org/zap v1.27.0 // indirect
	golang.org/x/exp v0.0.0-20240719175910-8a7402abbf56 // indirect
	golang.org/x/sync v0.12.0 // indirect
	gomodules.xyz/jsonpatch/v2 v2.5.0 // indirect
	google.golang.org/genproto/googleapis/api v0.0.0-20241209162323-e6fa225c2576 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20241223144023-3abc09e42ca8 // indirect
	gopkg.in/evanphx/json-patch.v4 v4.12.0 // indirect
	gopkg.in/go-jose/go-jose.v2 v2.6.3 // indirect
	k8s.io/klog/v2 v2.130.1 // indirect
	k8s.io/kube-openapi v0.31.0 // indirect
	k8s.io/utils v0.0.0-20241210054802-24370beab758 // indirect
	knative.dev/pkg v0.0.0-20231023150739-56bfe0dd9626 // indirect
	sigs.k8s.io/json v0.0.0-20241014173422-cfa47c3a1cc8 // indirect
	sigs.k8s.io/randfill v1.0.0 // indirect
	sigs.k8s.io/structured-merge-diff/v4 v4.6.0 // indirect
	sigs.k8s.io/yaml v1.4.0 // indirect
)

require (
	kubevirt.io/api v1.7.0
	kubevirt.io/client-go v1.7.0
	kubevirt.io/containerized-data-importer-api v1.63.1 // indirect
	kubevirt.io/controller-lifecycle-operator-sdk/api v0.0.0-20220329064328-f3cc58c6ed90 // indirect
)

replace (
	github.com/kubernetes-incubator/external-storage => github.com/libopenstorage/external-storage v0.20.4-openstorage-rc10
	github.com/libopenstorage/autopilot-api => github.com/libopenstorage/autopilot-api v0.6.1-0.20210301232050-ca2633c6e114
	github.com/portworx/torpedo => github.com/portworx/torpedo v0.0.0-20230206190621-4ccdccff9ded
	helm.sh/helm/v3 => helm.sh/helm/v3 v3.10.3
	// Override k8s.io/client-go incompatible version and fix invalid v0.0.0 versions from legacy dependencies
	// (external-storage and libopenstorage/operator have these as v0.0.0 placeholders)
	k8s.io/api => k8s.io/api v0.33.5
	k8s.io/apiextensions-apiserver => k8s.io/apiextensions-apiserver v0.33.5
	k8s.io/apimachinery => k8s.io/apimachinery v0.33.5
	k8s.io/apiserver => k8s.io/apiserver v0.33.5
	k8s.io/cli-runtime => k8s.io/cli-runtime v0.33.5
	k8s.io/client-go => k8s.io/client-go v0.33.5
	k8s.io/cloud-provider => k8s.io/cloud-provider v0.33.5
	k8s.io/cluster-bootstrap => k8s.io/cluster-bootstrap v0.33.5
	k8s.io/code-generator => k8s.io/code-generator v0.33.5
	k8s.io/component-base => k8s.io/component-base v0.33.5
	k8s.io/component-helpers => k8s.io/component-helpers v0.33.5
	k8s.io/controller-manager => k8s.io/controller-manager v0.33.5
	k8s.io/cri-api => k8s.io/cri-api v0.33.5
	k8s.io/cri-client => k8s.io/cri-client v0.33.5
	k8s.io/csi-translation-lib => k8s.io/csi-translation-lib v0.33.5
	k8s.io/dynamic-resource-allocation => k8s.io/dynamic-resource-allocation v0.33.5
	k8s.io/endpointslice => k8s.io/endpointslice v0.33.5
	k8s.io/externaljwt => k8s.io/externaljwt v0.33.5
	k8s.io/kms => k8s.io/kms v0.33.5
	k8s.io/kube-aggregator => k8s.io/kube-aggregator v0.33.5
	k8s.io/kube-controller-manager => k8s.io/kube-controller-manager v0.33.5
	k8s.io/kube-openapi => k8s.io/kube-openapi v0.0.0-20250318190949-c8a335a9a2ff
	k8s.io/kube-proxy => k8s.io/kube-proxy v0.33.5
	k8s.io/kube-scheduler => k8s.io/kube-scheduler v0.33.5
	k8s.io/kubectl => k8s.io/kubectl v0.33.5
	k8s.io/kubelet => k8s.io/kubelet v0.33.5
	k8s.io/kubernetes => k8s.io/kubernetes v1.33.5
	k8s.io/legacy-cloud-providers => k8s.io/legacy-cloud-providers v0.33.5
	k8s.io/metrics => k8s.io/metrics v0.33.5
	k8s.io/mount-utils => k8s.io/mount-utils v0.33.5
	k8s.io/pod-security-admission => k8s.io/pod-security-admission v0.33.5
	k8s.io/sample-apiserver => k8s.io/sample-apiserver v0.33.5
	kubevirt.io/containerized-data-importer-api v1.57.0-alpha1 => kubevirt.io/containerized-data-importer-api v1.56.1
)
