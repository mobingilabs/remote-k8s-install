package service

import (
	"fmt"
	"mobingi/ocean/pkg/tools/machine"
	"path/filepath"

	"mobingi/ocean/pkg/config"
	"mobingi/ocean/pkg/constants"
	cmdutil "mobingi/ocean/pkg/util/cmd"
	templateutil "mobingi/ocean/pkg/util/template"
)

const apiserverServiceTemplate = `[Unit]
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
type apiserverTemplateData struct {
	PrivateIP   string
	EtcdServers string
}

func newAPIServerTemplateData(cfg *config.Config) *apiserverTemplateData {
	//TODO util func can join more url together
	etcdServers := fmt.Sprintf("http://%s:2379", cfg.Masters[0].PrivateIP)
	return &apiserverTemplateData{
		EtcdServers: etcdServers,
		PrivateIP:   cfg.Masters[0].PrivateIP,
	}
}

func NewRunAPIServerJob(cfg *config.Config) (*machine.Job, error) {
	serviceData, err := getAPIServerServiceFile(cfg)
	if err != nil {
		return nil, err
	}

	job := machine.NewJob("kube-apiserver-service")

	job.AddCmd(cmdutil.NewWriteCmd(filepath.Join(constants.ServiceDir, constants.KubeApiserverService), string(serviceData)))
	job.AddCmd(cmdutil.NewSystemStartCmd(constants.KubeApiserverService))

	return job, nil
}

func getAPIServerServiceFile(cfg *config.Config) ([]byte, error) {
	data, err := templateutil.Parse(apiserverServiceTemplate, newAPIServerTemplateData(cfg))
	if err != nil {
		return nil, err
	}

	return data, nil
}