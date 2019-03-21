package dependence

import (
	"mobingi/ocean/pkg/tools/machine"
	cmdutil "mobingi/ocean/pkg/util/cmd"
)

func getSetEnvCommands() machine.CommandList {
	cl := machine.CommandList{}

	swapOffCmd := "swapoff -a"
	swapOffCheck := func(output string) bool {
		return true
	}
	cl.Add(swapOffCmd, swapOffCheck)

	k8sConfContent := "net.ipv4.ip_forward = 1"
	writeCmd := cmdutil.NewWriteCmd("/etc/sysctl.d/k8s.conf", k8sConfContent)
	writeCheck := func(output string) bool {
		return true
	}
	cl.Add(writeCmd, writeCheck)

	return cl
}
