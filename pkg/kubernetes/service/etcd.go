package service

import (
	"fmt"
	"path/filepath"
	"strconv"
	"strings"

	"mobingi/ocean/pkg/constants"
	"mobingi/ocean/pkg/tools/machine"
	cmdutil "mobingi/ocean/pkg/util/cmd"
	templateutil "mobingi/ocean/pkg/util/template"
)

const (
	EtcdPeerPort   = 2380
	EtcdClientPort = 2379
)

const etcdServiceTemplate = `[Unit]
Description=Etcd Server
After=network.target
Documentation=https://github.com/coreos

[Service]
ExecStart={{.BinDir}}/etcd \\
  --name={{.Name}} \\
	--initial-advertise-peer-urls=https://{{.IP}}:{{.PeerPort}} \\
  --listen-peer-urls=https://{{.IP}}:{{.PeerPort}} \\
  --listen-client-urls=https://{{.IP}}:{{.ClientPort}} \\
  --advertise-client-urls=https://{{.IP}}:{{.ClientPort}} \\
	--initial-cluster-token=etcd-cluster \\
	--initial-cluster={{.InitialCluster}} \\
  --initial-cluster-state=new \\
	--data-dir={{.DataDir}} \\
	--cert-file={{.PKIDir}}/server.crt \\
	--key-file={{.PKIDir}}/server.key \\
	--client-cert-auth=true \\
	--trusted-ca-file={{.PKIDir}}/ca.crt \\
	--peer-cert-file={{.PKIDir}}/peer.crt \\
	--peer-key-file={{.PKIDir}}/peer.key \\
	--peer-client-cert-auth=true \\
	--peer-trusted-ca-file={{.PKIDir}}/ca.crt \\
	--peer-key-file={{.PKIDir}}/peer.key
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
	PeerPort       string
	ClientPort     string

	PKIDir  string
	DataDir string
	BinDir  string
}

func newEtcdTemplateData(i int, ip, initialCluster, peerPort, clientPort string) *etcdTemplateData {
	templateData := &etcdTemplateData{
		Name:           getEtcdName(i),
		IP:             ip,
		InitialCluster: initialCluster,
	}

	if len(peerPort) == 0 {
		templateData.PeerPort = strconv.Itoa(EtcdPeerPort)
	} else {
		templateData.PeerPort = peerPort
	}

	if len(clientPort) == 0 {
		templateData.ClientPort = strconv.Itoa(EtcdClientPort)
	} else {
		templateData.ClientPort = clientPort
	}

	templateData.PKIDir = filepath.Join(constants.PKIDir, "etcd")
	templateData.DataDir = constants.ETCDDataDir
	templateData.BinDir = constants.BinDir

	return templateData
}

// NewRunEtcdJobs will write service file to disk and start etcd service
func NewRunEtcdJobs(clusterIPs []string, cretList map[string][]byte) ([]*machine.Job, error) {
	jobs := make([]*machine.Job, 0, len(clusterIPs))
	initalCluster := getInitialCluster(clusterIPs)
	for i, v := range clusterIPs {
		j := machine.NewJob("etcd-service")
		serviceData, err := getEtcdServiceFile(i, v, initalCluster)
		if err != nil {
			return nil, err
		}

		j.AddCmd(cmdutil.NewMkdirAllCmd(constants.WorkDir))
		j.AddCmd(cmdutil.NewMkdirAllCmd(constants.PKIDir))
		j.AddCmd(cmdutil.NewMkdirAllCmd(filepath.Join(constants.PKIDir, "etcd")))
		for k, v := range cretList {
			j.AddCmd(cmdutil.NewWriteCmd(filepath.Join(constants.WorkDir, "pki", k), string(v)))
		}

		j.AddCmd(cmdutil.NewWriteCmd(filepath.Join(constants.ServiceDir, constants.EtcdService), serviceData))
		j.AddCmd(cmdutil.NewSystemStartCmd(constants.EtcdService))
		jobs = append(jobs, j)
	}

	return jobs, nil
}

func getEtcdServiceFile(i int, ip, initialCluster string) (string, error) {
	data, err := templateutil.Parse(etcdServiceTemplate, newEtcdTemplateData(i, ip, initialCluster, "", ""))
	if err != nil {
		return "", err
	}

	return string(data), nil
}

func getInitialCluster(clusterIPs []string) string {
	urls := make([]string, 0, len(clusterIPs))
	for i, v := range clusterIPs {
		url := fmt.Sprintf("%s=https://%s:%d", getEtcdName(i), v, EtcdPeerPort)
		urls = append(urls, url)
	}

	return strings.Join(urls, ",")
}

func getEtcdName(num int) string {
	return fmt.Sprintf("etcd%d", num)
}

func GetEtcdServers(ips []string) string {
	urls := make([]string, 0, len(ips))
	for _, v := range ips {
		url := fmt.Sprintf("https://%s:%d", v, EtcdClientPort)
		urls = append(urls, url)
	}

	return strings.Join(urls, ",")
}
