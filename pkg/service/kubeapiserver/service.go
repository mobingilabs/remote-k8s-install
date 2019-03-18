package kubeapiserver

import (
	"mobingi/ocean/pkg/config"
)

const serviceTemplate = `[Unit]
Description=Kubernetes API Server
Documentation=https://github.com/GoogleCloudPlatform/kubernetes
After=network.target
After=etcd.service

[Service]
ExecStart=/usr/local/bin/kube-apiserver \
  --authorization-mode=Node,RBAC \
  --advertise-address={{.IP}} \
  --bind-address={{.IP}} \
  --client-ca-file=/etc/kubernetes/pki/ca.crt \
  --enable-admission-plugins=NodeRestriction \
  --enable-bootstrap-token-auth=true \
  --etcd-servers=http://{{.IP}}:2379 \
  --kubelet-client-certificate=/etc/kubernetes/pki/apiserver-kubelet-client.crt \
  --kubelet-client-key=/etc/kubernetes/pki/apiserver-kubelet-client.key \
  --kubelet-preferred-address-types=InternalIP,ExternalIP,Hostname \
  --proxy-client-cert-file=/etc/kubernetes/pki/front-proxy-client.crt \
  --proxy-client-key-file=/etc/kubernetes/pki/front-proxy-client.key \
  --requestheader-allowed-names=front-proxy-client \
  --requestheader-client-ca-file=/etc/kubernetes/pki/front-proxy-ca.crt \
  --requestheader-extra-headers-prefix=X-Remote-Extra- \
  --requestheader-group-headers=X-Remote-Group \
  --requestheader-username-headers=X-Remote-User \
  --secure-port=6443 \
  --service-account-key-file=/etc/kubernetes/pki/sa.pub \
  --tls-cert-file=/etc/kubernetes/pki/apiserver.crt \
  --tls-private-key-file=/etc/kubernetes/pki/apiserver.key \
Restart=on-failure
RestartSec=5

[Install]
WantedBy=multi-user.target`

// TODO more data from config to flexible
type templateData struct {
	IP string
}

func newTemplateData(cfg *config.Config) *templateData {
	masterMachine := cfg.GetMasterMachine()
	return &templateData{
		IP: masterMachine.Addr,
	}
}
