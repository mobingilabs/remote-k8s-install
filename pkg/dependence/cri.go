package dependence

import (
	"path/filepath"

	"mobingi/ocean/pkg/tools/machine"
	cmduitl "mobingi/ocean/pkg/util/cmd"
)

// SetupCRI install cri on the machine, now we use docker
func getCRICommands() machine.CommandList {
	cl := machine.CommandList{}

	yumCmd1 := "yum install -y yum-utils device-mapper-perisitent-data lvm2"
	// TODO real check
	yumCmdCheck1 := func(output string) bool {
		return true
	}
	cl.Add(yumCmd1, yumCmdCheck1)

	yumCmd2 := "yum-config-manager --add-repo https://download.docker.com/linux/centos/docker-ce.repo"
	yumCmdCheck2 := func(output string) bool {
		return true
	}
	cl.Add(yumCmd2, yumCmdCheck2)

	yumCmd3 := "yum install -y docker-ce-18.06.0.ce-3.el7"
	yumCmdCheck3 := func(output string) bool {
		return true
	}
	cl.Add(yumCmd3, yumCmdCheck3)

	// TODO dir we can use constans one location to lisa all
	dockerEtcDir := "/etc/docker"
	mkdirCmd := cmduitl.NewMkdirAllCmd(dockerEtcDir)
	mkdirCheck := func(output string) bool {
		return true
	}
	cl.Add(mkdirCmd, mkdirCheck)

	dockerDaemonJson := `{
		"exec-opts": ["native.cgroupdriver=systemd"],
		"log-driver": "json-file",
		"log-opts": {
			"max-size": "100m"
		},
		"storage-driver": "overlay2"
	}`
	dockerDaemonJsonName := "daemon.json"
	writeCmd := cmduitl.NewWriteCmd(filepath.Join(dockerEtcDir, dockerDaemonJsonName), dockerDaemonJson)
	writeCheck := func(output string) bool {
		return true
	}
	cl.Add(writeCmd, writeCheck)

	startCmd := cmduitl.NewSystemStartCmd("docker")
	startCheck := func(output string) bool {
		return true
	}
	cl.Add(startCmd, startCheck)

	return cl
}
