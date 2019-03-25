package node

import (
	"mobingi/ocean/pkg/config"
	"mobingi/ocean/pkg/constants"
	"mobingi/ocean/pkg/tools/machine"
	cmdutil "mobingi/ocean/pkg/util/cmd"
	checkutil "mobingi/ocean/pkg/util/check"
	"path/filepath"
)

func getWriteBootstrapCommands(data []byte) machine.CommandList {
	cl := machine.CommandList{}
	writeCmd := cmdutil.NewWriteCmd(filepath.Join(constants.WorkDir, constants.BootstrapKubeletConfName), string(data))
	cl.Add(writeCmd, checkutil.WriteCheck)
	return cl
}

func getDownloadCommands(cfg *config.Config) machine.CommandList {
	cl := machine.CommandList{}

	//TODO name get from constatns, writeCheck change to curl check
	cl.Add(cmdutil.NewCurlCmd(cfg.DownloadBinSite, "node.tgz"), checkutil.WriteCheck)
	cl.Add(cmdutil.NewCurlCmd(cfg.DownloadBinSite, "cni.tgz"), checkutil.WriteCheck)

	return cl
}
