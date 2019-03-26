package master

import (
	"mobingi/ocean/pkg/config"
	"path/filepath"

	"mobingi/ocean/pkg/constants"
	"mobingi/ocean/pkg/tools/machine"
	checkutil "mobingi/ocean/pkg/util/check"
	cmdutil "mobingi/ocean/pkg/util/cmd"
)

const (
	masterTgzName = "master.tgz"
	nodeTgzName   = "node.tgz"
	cniTgzName    = "cni.tgz"
)

func getWriteCertsCommand(certList map[string][]byte) machine.CommandList {
	cl := machine.CommandList{}

	for k, v := range certList {
		cmd := cmdutil.NewWriteCmd(filepath.Join(constants.PKIDir, k), string(v))
		cl.Add(cmd, checkutil.WriteCheck)
	}

	return cl
}

func getWriteKubeconfsCommand(kubeconfs map[string][]byte) machine.CommandList {
	cl := machine.CommandList{}

	for k, v := range kubeconfs {
		cmd := cmdutil.NewWriteCmd(filepath.Join(constants.WorkDir, k), string(v))
		cl.Add(cmd, checkutil.WriteCheck)
	}

	return cl
}

func getDownloadCommands(cfg *config.Config) machine.CommandList {
	cl := machine.CommandList{}

	//TODO name get from constatns, writeCheck change to curl check
	cl.Add(cmdutil.NewCurlCmd(cfg.DownloadBinSite, "master.tgz"), checkutil.WriteCheck)

	return cl
}