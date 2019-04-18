package service

import (
	"path/filepath"

	"mobingi/ocean/pkg/constants"
	"mobingi/ocean/pkg/tools/machine"
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
	--allow-privileged=true \\
  --advertise-address={{.IP}} \\
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
	--tls-private-key-file=/etc/kubernetes/pki/apiserver.key \\
	--etcd-cafile=/etc/kubernetes/pki/etcd/ca.crt \\
	--etcd-certfile=/etc/kubernetes/pki/apiserver-etcd-client.crt \\
	--etcd-keyfile=/etc/kubernetes/pki/apiserver-etcd-client.key 
Restart=on-failure
RestartSec=5

[Install]
WantedBy=multi-user.target`

// TODO more data from config to flexible
type apiserverTemplateData struct {
	// now it is a private ip
	IP               string
	EtcdServers      string
	AdvertiseAddress string
}

func NewRunAPIServerJobs(ips []string, etcdServers, advertiseAddress string) ([]*machine.Job, error) {
	jobs := make([]*machine.Job, 0, len(ips))
	for _, v := range ips {
		serviceData, err := getAPIServerServiceFile(v, etcdServers, advertiseAddress)
		if err != nil {
			return nil, err
		}

		job := machine.NewJob("kube-apiserver-service")
		job.AddCmd(cmdutil.NewWriteCmd(filepath.Join(constants.ServiceDir, constants.KubeApiserverService), string(serviceData)))
		job.AddCmd(cmdutil.NewSystemStartCmd(constants.KubeApiserverService))
		jobs = append(jobs, job)
	}

	return jobs, nil
}

func getAPIServerServiceFile(ip, etcdServers, advertiseAddress string) ([]byte, error) {
	templateData := apiserverTemplateData{
		IP:               ip,
		EtcdServers:      etcdServers,
		AdvertiseAddress: advertiseAddress,
	}
	data, err := templateutil.Parse(apiserverServiceTemplate, templateData)
	if err != nil {
		return nil, err
	}

	return data, nil
}
