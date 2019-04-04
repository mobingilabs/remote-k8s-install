package node

import (
	"mobingi/ocean/pkg/constants"
	"mobingi/ocean/pkg/tools/machine"
	cmdutil "mobingi/ocean/pkg/util/cmd"
	"path/filepath"
)

func NewJob(site string, bootstrapconf []byte) *machine.Job {
	j := machine.NewJob("node-prepare")

	mkdir(j)
	setEnv(j)
	writeBootstrapconf(j, bootstrapconf)
	downloadBin(j, site)
	installDocker(j)

	return j
}

func mkdir(j *machine.Job) {
	j.AddCmd(cmdutil.NewMkdirAllCmd(constants.WorkDir))
	j.AddCmd(cmdutil.NewMkdirAllCmd(constants.PKIDir))
	j.AddCmd(cmdutil.NewMkdirAllCmd("/opt/bin/cni"))
}

func setEnv(j *machine.Job) {
	j.AddCmd("swapoff -a")
	j.AddCmd(cmdutil.NewWriteCmd("/etc/sysctl.d/k8s.conf", "net.ipv4.ip_forward = 1"))
}

func writeBootstrapconf(j *machine.Job, bootstrapconf []byte) {
	j.AddCmd(cmdutil.NewWriteCmd(filepath.Join(constants.WorkDir, constants.BootstrapKubeletConfName), string(bootstrapconf)))
}

func downloadBin(j *machine.Job, site string) {
	j.AddCmd(cmdutil.NewCurlCmd(site, "master.tgz"))
	j.AddCmd(cmdutil.NewCurlCmd(site, "cni.tgz"))
	j.AddCmd(cmdutil.NewTarXCmd("/tmp/master.tgz", "/usr/local/bin"))
	j.AddCmd(cmdutil.NewTarXCmd("/tmp/cni.tgz", "/opt/bin/cni"))
}

func installDocker(j *machine.Job) {
	j.AddCmd("yum install -y yum-utils device-mapper-perisitent-data lvm2")
	j.AddCmd("yum-config-manager --add-repo https://download.docker.com/linux/centos/docker-ce.repo")
	j.AddCmd("yum install -y docker-ce-18.06.0.ce-3.el7")
	j.AddCmd(cmdutil.NewMkdirAllCmd("/etc/docker"))

	dockerDaemonJson := `{
	"exec-opts": ["native.cgroupdriver=systemd"],
	"log-driver": "json-file",
	"log-opts": {
		"max-size": "100m"
	},
	"storage-driver": "overlay2"
}`
	dockerDaemonJsonName := "daemon.json"
	j.AddCmd(cmdutil.NewWriteCmd(filepath.Join("/etc/docker", dockerDaemonJsonName), dockerDaemonJson))
	// TODO add check
	j.AddCmd(cmdutil.NewSystemStartCmd("docker"))
}
