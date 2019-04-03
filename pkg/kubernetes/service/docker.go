package service

import (
	"mobingi/ocean/pkg/config"
	"mobingi/ocean/pkg/tools/machine"
	cmdutil "mobingi/ocean/pkg/util/cmd"
)

func NewRunDockerJob(cfg *config.Config) *machine.Job {
	job := machine.NewJob("docker-service")

	job.AddCmd("yum install -y yum-utils device-mapper-persistent-data lvm2")
	job.AddCmd("yum-config-manager --add-repo  https://download.docker.com/linux/centos/docker-ce.repo")
	job.AddCmd("yum install -y docker-ce-18.06.0.ce-3.el7")
	// may be node make dir do this
	job.AddCmd(cmdutil.NewMkdirAllCmd("/etc/docker"))

	writeContent := `{
	"exec-opts": ["native.cgroupdriver=systemd"],
	"log-driver": "json-file",
	"log-opts": {
		"max-size": "100m"
	},
	"storage-driver": "overlay2"
}`
	job.AddCmd(cmdutil.NewWriteCmd("/etc/docker/daemon.json", writeContent))

	job.AddCmd(cmdutil.NewSystemStartCmd("docker.service"))

	return job
}
