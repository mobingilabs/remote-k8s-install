package phases

import (
	"mobingi/ocean/pkg/tools/machine"
	cmdutil "mobingi/ocean/pkg/util/cmd"
	"path/filepath"
)

func InstallDocker(j *machine.Job) {
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
