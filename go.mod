module github.com/portworx/sched-ops

go 1.22.6

require (
	github.com/coreos/prometheus-operator v0.31.1
	github.com/golang/mock v1.2.0
	github.com/kubernetes-csi/external-snapshotter/v2 v2.1.1
	github.com/kubernetes-incubator/external-storage v0.0.0-00010101000000-000000000000
	github.com/libopenstorage/autopilot-api v0.6.0
	github.com/libopenstorage/openstorage v8.0.0+incompatible
	github.com/libopenstorage/operator v0.0.0-20200725001727-48d03e197117
	github.com/libopenstorage/stork v1.3.0-beta1.0.20200630005842-9255e7a98775
	github.com/openshift/api v0.0.0-20190322043348-8741ff068a47
	github.com/openshift/client-go v0.0.0-20180830153425-431ec9a26e50
	github.com/portworx/talisman v0.0.0-20191007232806-837747f38224
	github.com/sirupsen/logrus v1.4.2
	github.com/stretchr/testify v1.4.0
	github.com/undefinedlabs/go-mpatch v1.0.7
	k8s.io/api v0.17.0
	k8s.io/apiextensions-apiserver v0.0.0
	k8s.io/apimachinery v0.17.1-beta.0
	k8s.io/client-go v11.0.1-0.20190409021438-1a26190bd76a+incompatible
)

require (
	github.com/PuerkitoBio/purell v1.1.1 // indirect
	github.com/PuerkitoBio/urlesc v0.0.0-20170810143723-de5bf2ad4578 // indirect
	github.com/coreos/go-oidc v0.0.0-20180117170138-065b426bd416 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/dgrijalva/jwt-go v3.2.0+incompatible // indirect
	github.com/docker/spdystream v0.0.0-20181023171402-6480d4af844c // indirect
	github.com/elazarl/goproxy v0.0.0-20191011121108-aa519ddbe484 // indirect
	github.com/emicklei/go-restful v2.10.0+incompatible // indirect
	github.com/go-openapi/jsonpointer v0.19.3 // indirect
	github.com/go-openapi/jsonreference v0.19.3 // indirect
	github.com/go-openapi/spec v0.19.3 // indirect
	github.com/go-openapi/swag v0.19.5 // indirect
	github.com/gogo/protobuf v1.3.0 // indirect
	github.com/golang/glog v0.0.0-20160126235308-23def4e6c14b // indirect
	github.com/golang/groupcache v0.0.0-20190129154638-5b532d6fd5ef // indirect
	github.com/golang/protobuf v1.3.2 // indirect
	github.com/google/gofuzz v1.0.0 // indirect
	github.com/googleapis/gnostic v0.3.1 // indirect
	github.com/grpc-ecosystem/grpc-gateway v1.11.3 // indirect
	github.com/imdario/mergo v0.3.8 // indirect
	github.com/json-iterator/go v1.1.7 // indirect
	github.com/konsorten/go-windows-terminal-sequences v1.0.2 // indirect
	github.com/mailru/easyjson v0.7.0 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.1 // indirect
	github.com/mohae/deepcopy v0.0.0-20170929034955-c48cc78d4826 // indirect
	github.com/operator-framework/operator-sdk v0.0.7 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/pquerna/cachecontrol v0.0.0-20180517163645-1555304b9b35 // indirect
	github.com/spf13/pflag v1.0.3 // indirect
	golang.org/x/crypto v0.0.0-20191010185427-af544f31c8ac // indirect
	golang.org/x/net v0.0.0-20191209160850-c0dbc17a3553 // indirect
	golang.org/x/oauth2 v0.0.0-20190604053449-0f29369cfe45 // indirect
	golang.org/x/sys v0.0.0-20191220220014-0732a990476f // indirect
	golang.org/x/text v0.3.2 // indirect
	golang.org/x/time v0.0.0-20190921001708-c4c64cad1fd0 // indirect
	google.golang.org/appengine v1.6.5 // indirect
	google.golang.org/genproto v0.0.0-20191220175831-5c49e3ecc1c1 // indirect
	google.golang.org/grpc v1.26.0 // indirect
	gopkg.in/inf.v0 v0.9.1 // indirect
	gopkg.in/square/go-jose.v2 v2.3.1 // indirect
	gopkg.in/yaml.v2 v2.2.8 // indirect
	k8s.io/klog v1.0.0 // indirect
	k8s.io/kube-openapi v0.0.0-20190918143330-0270cf2f1c1d // indirect
	k8s.io/utils v0.0.0-20190923111123-69764acb6e8e // indirect
	sigs.k8s.io/controller-runtime v0.2.2 // indirect
	sigs.k8s.io/yaml v1.1.0 // indirect
)

replace (
	github.com/kubernetes-incubator/external-storage => github.com/libopenstorage/external-storage v5.1.1-0.20190919185747-9394ee8dd536+incompatible
	github.com/kubernetes-incubator/external-storage v0.0.0-00010101000000-000000000000 => github.com/libopenstorage/external-storage v5.1.1-0.20190919185747-9394ee8dd536+incompatible
	k8s.io/api => k8s.io/api v0.15.11
	k8s.io/apiextensions-apiserver => k8s.io/apiextensions-apiserver v0.15.11
	k8s.io/apimachinery => k8s.io/apimachinery v0.15.11
	k8s.io/apiserver => k8s.io/apiserver v0.15.11
	k8s.io/cli-runtime => k8s.io/cli-runtime v0.15.11
	k8s.io/client-go => k8s.io/client-go v0.15.11
	k8s.io/cloud-provider => k8s.io/cloud-provider v0.15.11
	k8s.io/cluster-bootstrap => k8s.io/cluster-bootstrap v0.15.11
	k8s.io/code-generator => k8s.io/code-generator v0.15.11
	k8s.io/component-base => k8s.io/component-base v0.15.11
	k8s.io/cri-api => k8s.io/cri-api v0.15.11
	k8s.io/csi-translation-lib => k8s.io/csi-translation-lib v0.15.11
	k8s.io/kube-aggregator => k8s.io/kube-aggregator v0.15.11
	k8s.io/kube-controller-manager => k8s.io/kube-controller-manager v0.15.11
	k8s.io/kube-proxy => k8s.io/kube-proxy v0.15.11
	k8s.io/kube-scheduler => k8s.io/kube-scheduler v0.15.11
	k8s.io/kubectl => k8s.io/kubectl v0.15.11
	k8s.io/kubelet => k8s.io/kubelet v0.15.11
	k8s.io/kubernetes => k8s.io/kubernetes v1.16.0
	k8s.io/legacy-cloud-providers => k8s.io/legacy-cloud-providers v0.15.11
	k8s.io/metrics => k8s.io/metrics v0.15.11
	k8s.io/sample-apiserver => k8s.io/sample-apiserver v0.15.11
)
