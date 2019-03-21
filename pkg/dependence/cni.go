package dependence

import (
	"mobingi/ocean/pkg/constants"
	"mobingi/ocean/pkg/tools/machine"
	cmdutil "mobingi/ocean/pkg/util/cmd"
)

// kubelet need cni plugin for network
func getCniBinCommands() machine.CommandList {
	cl := machine.CommandList{}

	// TODO from config
	curlCmd := cmdutil.NewCurlCmd(targetSite, cniTgzName)
	curlCheck := func(output string) bool {
		return true
	}
	cl.Add(curlCmd, curlCheck)

	tarCmd := cmdutil.NewTarXCmd(cniTgzName, constants.BinDir)
	tarCheck := func(output string) bool {
		return true
	}
	cl.Add(tarCmd, tarCheck)

	return cl
}

