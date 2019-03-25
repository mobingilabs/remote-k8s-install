package dependence

import (
	"mobingi/ocean/pkg/config"
	"mobingi/ocean/pkg/constants"
	"mobingi/ocean/pkg/tools/machine"
	cmdutil "mobingi/ocean/pkg/util/cmd"
)

// get from fileServer
const (
	masterTgzName = "master.tgz"
	nodeTgzName   = "node.tgz"
	cniTgzName    = "cni.tgz"
	targetSite    = "localhost:8080/"
)

func getMasterBinCommands() machine.CommandList {
	cl := machine.CommandList{}

	// TODO from config
	curlCmd := cmdutil.NewCurlCmd(targetSite, masterTgzName)
	curlCheck := func(output string) bool {
		return true
	}
	cl.Add(curlCmd, curlCheck)

	tarCmd := cmdutil.NewTarXCmd(masterTgzName, constants.BinDir)
	tarCheck := func(output string) bool {
		return true
	}
	cl.Add(tarCmd, tarCheck)

	return cl
}

func getNodeBinCommands(cfg *config.Config) machine.CommandList {
	cl := machine.CommandList{}

	// TODO from config
	curlCmd := cmdutil.NewCurlCmd(cfg.DownloadBinSite, nodeTgzName)
	curlCheck := func(output string) bool {
		return true
	}
	cl.Add(curlCmd, curlCheck)

	tarCmd := cmdutil.NewTarXCmd(nodeTgzName, constants.BinDir)
	tarCheck := func(output string) bool {
		return true
	}
	cl.Add(tarCmd, tarCheck)

	return cl
}
