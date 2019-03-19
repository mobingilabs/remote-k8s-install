package kubeapiserver

import (
	"fmt"
	"mobingi/ocean/pkg/config"
)

// TODO more info from config and constants to centralized
const serviceTemplate = `[Unit]
Description=Kubernetes API Server
Documentation=https://github.com/GoogleCloudPlatform/kubernetes
After=network.target
After=etcd.service

[Service]
ExecStart=/usr/local/bin/kube-apiserver \\
  --authorization-mode=Node,RBAC \\
  --advertise-address={{.PrivateIP}} \\
  --bind-address=0.0.0.0 \\
  --client-ca-file=/etc/kubernetes/pki/ca.crt \\
  --enable-admission-plugins=NodeRestriction \\
  --enable-bootstrap-token-auth=true \\
  --etcd-servers={{.EtcdServers}} \\
  --kubelet-client-certificate=/etc/kubernetes/pki/apiserver-kubelet-client.crt \\
  --kubelet-client-key=/etc/kubernetes/pki/apiserver-kubelet-client.key \\
  --kubelet-preferred-address-types=InternalIP,ExternalIP,Hostname \\
  --proxy-client-cert-file=/etc/kubernetes/pki/front-proxy-client.crt \\
  --proxy-client-key-file=/etc/kubernetes/pki/front-proxy-client.key \\
  --requestheader-allowed-names=front-proxy-client \\
  --requestheader-client-ca-file=/etc/kubernetes/pki/front-proxy-ca.crt \\
  --requestheader-extra-headers-prefix=X-Remote-Extra- \\
  --requestheader-group-headers=X-Remote-Group \\
  --requestheader-username-headers=X-Remote-User \\
  --secure-port=6443 \\
  --service-account-key-file=/etc/kubernetes/pki/sa.pub \\
  --tls-cert-file=/etc/kubernetes/pki/apiserver.crt \\
  --tls-private-key-file=/etc/kubernetes/pki/apiserver.key 
Restart=on-failure
RestartSec=5

[Install]
WantedBy=multi-user.target`

// TODO more data from config to flexible
type templateData struct {
	PrivateIP   string
	EtcdServers string
}

func newTemplateData(cfg *config.Config) *templateData {
	//TODO util func can join more url together
	etcdServers := fmt.Sprintf("http://%s:2379", cfg.Masters[0].PrivateIP)
	return &templateData{
		EtcdServers: etcdServers,
		PrivateIP:   cfg.Masters[0].PrivateIP,
	}
}
