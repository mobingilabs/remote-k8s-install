package docker

import (
	"mobingi/ocean/pkg/config"
	"mobingi/ocean/pkg/tools/machine"
	checkutil "mobingi/ocean/pkg/util/check"
	cmdutil "mobingi/ocean/pkg/util/cmd"
)

func CommandList(cfg *config.Config) machine.CommandList {
	cl := machine.CommandList{}

	yum0 := "yum install -y yum-utils device-mapper-persistent-data lvm2"
	cl.Add(yum0, checkutil.YUMCheck)

	yum1 := "yum-config-manager --add-repo  https://download.docker.com/linux/centos/docker-ce.repo"
	cl.Add(yum1, checkutil.YUMCheck)

	yum2 := "yum install -y docker-ce-18.06.0.ce-3.el7"
	cl.Add(yum2, checkutil.YUMCheck)

	mkdirCmd := cmdutil.NewMkdirAllCmd("/etc/docker")
	cl.Add(mkdirCmd, checkutil.MkdirCheck)

	writeContent := `{
{
	"exec-opts": ["native.cgroupdriver=systemd"],
	"log-driver": "json-file",
	"log-opts": {
		"max-size": "100m"
	},
	"storage-driver": "overlay2"
}`
	writeCmd := cmdutil.NewWriteCmd("/etc/docker/daemon.json", writeContent)
	cl.Add(writeCmd, checkutil.WriteCheck)

	systemdStartCmd := cmdutil.NewSystemStartCmd("docker.service")
	cl.Add(systemdStartCmd, checkutil.SystemStartCheck)

	return cl
}
