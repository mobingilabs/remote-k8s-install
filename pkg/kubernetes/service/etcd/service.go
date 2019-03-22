package etcd

import (
	"mobingi/ocean/pkg/config"
)

const serviceTemplate = `[Unit]
Description=Etcd Server
After=network.target
Documentation=https://github.com/coreos

[Service]
ExecStart=/usr/local/bin/etcd \\
  --name=node0 \\
  --data-dir=/var/lib/etcd \\
  --listen-client-urls=http://{{.IP}}:2379 \\
  --advertise-client-urls=http://{{.IP}}:2379 \\
  --listen-peer-urls=http://{{.IP}}:2380 \\
  --initial-advertise-peer-urls=http://{{.IP}}:2380 \\
  --initial-cluster=node0=http://{{.IP}}:2380 \\
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

// etcd is just listen private net, so we just need private ip to service file
func newTemplateData(cfg *config.Config) *templateData {
	return &templateData{
		IP: cfg.Masters[0].PrivateIP,
	}
}
