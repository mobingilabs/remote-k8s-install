module mobingi/ocean

require (
	github.com/docker/distribution v2.7.1+incompatible // indirect
	github.com/evanphx/json-patch v4.1.0+incompatible // indirect
	github.com/gogo/protobuf v1.2.1 // indirect
	github.com/golang/protobuf v1.3.1 // indirect
	github.com/google/gofuzz v1.0.0 // indirect
	github.com/googleapis/gnostic v0.2.0 // indirect
	github.com/imdario/mergo v0.3.7 // indirect
	github.com/json-iterator/go v1.1.6 // indirect
	github.com/lithammer/dedent v1.1.0 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.1 // indirect
	github.com/onsi/ginkgo v1.8.0 // indirect
	github.com/onsi/gomega v1.5.0 // indirect
	github.com/opencontainers/go-digest v1.0.0-rc1 // indirect
	github.com/pkg/errors v0.8.1
	github.com/spf13/afero v1.2.2 // indirect
	github.com/spf13/pflag v1.0.3 // indirect
	github.com/stretchr/testify v1.3.0
	golang.org/x/crypto v0.0.0-20190325154230-a5d413f7728c
	golang.org/x/net v0.0.0-20190328230028-74de082e2cca // indirect
	golang.org/x/oauth2 v0.0.0-20190319182350-c85d3e98c914 // indirect
	golang.org/x/time v0.0.0-20190308202827-9d24e82272b4 // indirect
	gopkg.in/inf.v0 v0.9.1 // indirect
	gopkg.in/yaml.v2 v2.2.2
	k8s.io/api v0.0.0-20190416052506-9eb4726e83e4
	k8s.io/apiextensions-apiserver v0.0.0-20190330190201-4cac3cbacb4e // indirect
	k8s.io/apimachinery v0.0.0-20190416052411-7dcd37fca543
	k8s.io/apiserver v0.0.0-20190401145308-d20c276e0982 // indirect
	k8s.io/client-go v11.0.0+incompatible
	k8s.io/cloud-provider v0.0.0-20190323031113-9c9d72d1bf90 // indirect
	k8s.io/cluster-bootstrap v0.0.0-20190313124217-0fa624df11e9
	k8s.io/component-base v0.0.0-20190313120452-4727f38490bc // indirect
	k8s.io/klog v0.2.0 // indirect
	k8s.io/kube-openapi v0.0.0-20190401085232-94e1e7b7574c // indirect
	k8s.io/kube-proxy v0.0.0-20190320190624-78a1c9778e0e // indirect
	k8s.io/kubelet v0.0.0-20190313123811-3556bcde9670 // indirect
	k8s.io/kubernetes v1.14.1
	k8s.io/utils v0.0.0-20190308190857-21c4ce38f2a7 // indirect
	sigs.k8s.io/yaml v1.1.0 // indirect
)

replace (
	k8s.io/api => k8s.io/api v0.0.0-20190313235455-40a48860b5ab
	k8s.io/apimachinery => k8s.io/apimachinery v0.0.0-20190404173353-6a84e37a896d
	k8s.io/client-go => k8s.io/client-go v11.0.0+incompatible
)
