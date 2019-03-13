package etcd

import (
	"mobingi/ocean/pkg/config"
)

const serviceTemplate = `[Unit]
Description=Etcd Server
After=network.target
Documentation=https://github.com/coreos

[Service]
ExecStart=/usr/local/bin/etcd \
  --name=node0 \
  --data-dir=/var/lib/etcd \
  --listen-client-urls=https://{{.IP}}:2379 \
  --advertise-client-urls=https://{{.IP}}:2379 \
  --listen-peer-urls=https://{{.IP}}:2380 \
  --initial-advertise-peer-urls=https://{{.IP}}:2380 \
  --initial-cluster=node0=https://{{.IP}}:2380 \
  --trusted-ca-file=/etc/kubernetes/pki/etcd/ca.crt \
  --cert-file=/etc/kubernetes/pki/etcd/server.crt \
  --key-file=/etc/kubernetes/pki/etcd/server.key \
  --client-cert-auth=true \
  --peer-trusted-ca-file=/etc/kubernetes/pki/etcd/ca.crt \
  --peer-cert-file=/etc/kubernetes/pki/etcd/peer.crt \
  --peer-key-file=/etc/kubernetes/pki/etcd/peer.key \
  --peer-client-cert-auth=true \
  --initial-cluster-token=tkn \
  --initial-cluster-state=new
Type=notify
Restart=on-failure
RestartSec=5
LimitNOFILE=65536

[Install]
WantedBy=multi-user.target`

type templateData struct {
	IP string
}

func newTemplateData(cfg *config.Config) *templateData {
	masterMachine := cfg.GetMasterMachine()
	return &templateData{
		IP: masterMachine.Addr,
	}
}
