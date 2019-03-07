##关于手动搭建的问题总结
### 1证书
网上手动搭建教程
证书配置文件:ca-config.json
```json
{
  "signing": {
    "default": {
      "expiry": "87600h"
    },
    "profiles": {
      "kubernetes": {
        "usages": [
            "signing",
            "key encipherment",
            "server auth",
            "client auth"
        ],
        "expiry": "87600h"
      }
    }
  }
}
```
证书请求文件:ca-csr.json
```json
{
  "CN": "kubernetes",
  "key": {
    "algo": "rsa",
    "size": 2048
  },
  "names": [
    {
      "C": "TW",
      "ST": "Taipei",
      "L": "Taipei",
      "O": "Kubernetes",
      "OU": "Kubernetes-manual"
    }
  ]
}
```
证书和私钥生成文件
```shell
cfssl gencert \
  -ca=${PKI_DIR}/ca.pem \
  -ca-key=${PKI_DIR}/ca-key.pem \
  -config=ca-config.json \
  -profile=kubernetes \
  admin-csr.json | cfssljson -bare ${PKI_DIR}/admin
```
问题，profile单纯的用了一套，所有ext权限全都赋予子证书，CN表示特定的对象，在某些情况下k8s会用此作为用户名，来给予默认的权限,手动生成复杂,
hosts配置，难以给全,和特定的机器IP耦合
kubeadm中的证书配置
```go
KubeadmCertRootCA = KubeadmCert{
		Name:     "ca",
		LongName: "self-signed Kubernetes CA to provision identities for other Kubernetes components",
		BaseName: kubeadmconstants.CACertAndKeyBaseName,
		config: certutil.Config{
			CommonName: "kubernetes",
		},
	}
	// KubeadmCertAPIServer is the definition of the cert used to serve the Kubernetes API.
	KubeadmCertAPIServer = KubeadmCert{
		Name:     "apiserver",
		LongName: "certificate for serving the Kubernetes API",
		BaseName: kubeadmconstants.APIServerCertAndKeyBaseName,
		CAName:   "ca",
		config: certutil.Config{
			CommonName: kubeadmconstants.APIServerCertCommonName,
			Usages:     []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		},
		configMutators: []configMutatorsFunc{
			makeAltNamesMutator(pkiutil.GetAPIServerAltNames),
		},
	}
```
官方代码示例，configmutators来根据配置负责配置
host，config的权限也是根据需求给予
2.kubeconfig文件问题
master组件需要使用
教程
```shell
kubectl config set-cluster kubernetes \
    --certificate-authority=${PKI_DIR}/ca.pem \
    --embed-certs=true \
    --server=${KUBE_APISERVER} \
    --kubeconfig=${K8S_DIR}/scheduler.conf

$ kubectl config set-credentials system:kube-scheduler \
    --client-certificate=${PKI_DIR}/scheduler.pem \
    --client-key=${PKI_DIR}/scheduler-key.pem \
    --embed-certs=true \
    --kubeconfig=${K8S_DIR}/scheduler.conf

$ kubectl config set-context system:kube-scheduler@kubernetes \
    --cluster=kubernetes \
    --user=system:kube-scheduler \
    --kubeconfig=${K8S_DIR}/scheduler.conf

$ kubectl config use-context system:kube-scheduler@kubernetes \
    --kubeconfig=${K8S_DIR}/scheduler.conf
```
官方代码
```go
    config := CreateBasic(serverURL, clusterName, userName, caCert)
	config.AuthInfos[userName] = &clientcmdapi.AuthInfo{
		ClientKeyData:         clientKey,
		ClientCertificateData: clientCert,
	}
	return config

  type Config struct {
	// Legacy field from pkg/api/types.go TypeMeta.
	// TODO(jlowdermilk): remove this after eliminating downstream dependencies.
	// +optional
	Kind string `json:"kind,omitempty"`
	// Legacy field from pkg/api/types.go TypeMeta.
	// TODO(jlowdermilk): remove this after eliminating downstream dependencies.
	// +optional
	APIVersion string `json:"apiVersion,omitempty"`
	// Preferences holds general information to be use for cli interactions
	Preferences Preferences `json:"preferences"`
	// Clusters is a map of referencable names to cluster configs
	Clusters map[string]*Cluster `json:"clusters"`
	// AuthInfos is a map of referencable names to user configs
	AuthInfos map[string]*AuthInfo `json:"users"`
	// Contexts is a map of referencable names to context configs
	Contexts map[string]*Context `json:"contexts"`
	// CurrentContext is the name of the context that you would like to use by default
	CurrentContext string `json:"current-context"`
	// Extensions holds additional information. This is useful for extenders so that reads and writes don't clobber unknown fields
	// +optional
	Extensions map[string]runtime.Object `json:"extensions,omitempty"`
}
```
我用kubectl配置，发现遇到权限问题，也没找到问题所在，看代码
client生成的时候会用这个结构体的配置，支持多namespaces，context，user。方便切换,而且考虑到了多集群的切换问题