package service

import (
	"fmt"
	"path/filepath"
	"strings"

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
	--name={{.Name}} \\
	--initial-advertise-peer-urls=http://{{.IP}}:2380 \\
  --listen-peer-urls=http://{{.IP}}:2380 \\
  --listen-client-urls=http://{{.IP}}:2379 \\
  --advertise-client-urls=http://{{.IP}}:2379 \\
	--initial-cluster-token=etcd-cluster
	--initial-cluster={{.INITIAL_CLUSTER}}
	--initial-cluster-state=new \\
  --data-dir=/var/lib/etcd 
Type=notify
Restart=on-failure
RestartSec=5
LimitNOFILE=65536

[Install]
WantedBy=multi-user.target`

type etcdTemplateData struct {
	Name           string
	IP             string
	InitialCluster string
}

func newEtcdTemplateData(i int, ip, initialCluster string) *etcdTemplateData {
	return &etcdTemplateData{
		Name:           getEtcdName(i),
		IP:             ip,
		InitialCluster: initialCluster,
	}
}

func NewRunEtcdJobs(clusterIPs []string) ([]*machine.Job, error) {
	jobs := make([]*machine.Job, 0, len(clusterIPs))
	initalCluster := getInitalCluster(clusterIPs)
	for i, v := range clusterIPs {
		j := machine.NewJob("etcd-service")
		serviceData, err := getEtcdServiceFile(i, v, initalCluster)
		if err != nil {
			return nil, err
		}
		j.AddCmd(cmdutil.NewWriteCmd(filepath.Join(constants.ServiceDir, constants.EtcdService), string(serviceData)))
		j.AddCmd(cmdutil.NewSystemStartCmd(constants.EtcdService))
		jobs = append(jobs, j)
	}

	return jobs, nil
}

func getEtcdServiceFile(i int, ip, initialCluster string) ([]byte, error) {
	data, err := templateutil.Parse(etcdServiceTemplate, newEtcdTemplateData(i, ip, initialCluster))
	if err != nil {
		return nil, err
	}

	return data, nil
}

func getInitalCluster(clusterIPs []string) string {
	urls := make([]string, 0, len(clusterIPs))
	for i, v := range clusterIPs {
		url := fmt.Sprintf("%s=http://%s:2380", getEtcdName(i), v)
		urls = append(urls, url)
	}

	return strings.Join(urls, ",")
}

func getEtcdName(num int) string {
	return fmt.Sprintf("etcd%d", num)
}

func GetEtcdServers (ips []string) string {
	urls := make([]string, 0, len(ips))
	for _, v := range ips {
		url := fmt.Sprintf("http://%s:2379", v)
		urls = append(urls, url)
	}

	return strings.Join(urls, ",")
}

