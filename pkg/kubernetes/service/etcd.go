package service

import (
	"path/filepath"

	"mobingi/ocean/pkg/config"
	"mobingi/ocean/pkg/constants"
	"mobingi/ocean/pkg/tools/machine"
	cmdutil "mobingi/ocean/pkg/util/cmd"
	templateutil "mobingi/ocean/pkg/util/template"
)

const etcdServiceTemplate = `[Unit]
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

type etcdTemplateData struct {
	IP string
}

// etcd is just listen private net, so we just need private ip to service file
func newEtcdTemplateData(cfg *config.Config) *etcdTemplateData {
	return &etcdTemplateData{
		IP: cfg.Masters[0].PrivateIP,
	}
}

func NewRunEtcdJob(cfg *config.Config) (*machine.Job, error) {
	j := machine.NewJob("etcd-service")
	serviceData, err := getEtcdServiceFile(cfg)
	if err != nil {
		return nil, err
	}

	j.AddCmd(cmdutil.NewWriteCmd(filepath.Join(constants.ServiceDir, constants.EtcdService), string(serviceData)))
	j.AddCmd(cmdutil.NewSystemStartCmd(constants.EtcdService))

	return j, nil
}

func getEtcdServiceFile(cfg *config.Config) ([]byte, error) {
	data, err := templateutil.Parse(etcdServiceTemplate, newEtcdTemplateData(cfg))
	if err != nil {
		return nil, err
	}

	return data, nil
}
