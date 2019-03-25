package dependence

import (
	"mobingi/ocean/pkg/config"
	"mobingi/ocean/pkg/tools/machine"
	checkutil "mobingi/ocean/pkg/util/check"
	cmdutil "mobingi/ocean/pkg/util/cmd"
)

func GetMasterSetCommands() machine.CommandList {
	cl := machine.CommandList{}
	cl.AddAnother(getSetEnvCommands())
	cl.Add(cmdutil.NewTarXCmd("/tmp/master.tgz", "/usr/local/bin"), checkutil.SCPCheck)

	return cl
}

func GetNodesSetCommands(cfg *config.Config) machine.CommandList {
	cl := machine.CommandList{}
	cl.AddAnother(getSetEnvCommands())
	cl.Add(cmdutil.NewTarXCmd("/tmp/node.tgz", "/usr/local/bin"), checkutil.SCPCheck)
	cl.Add(cmdutil.NewTarXCmd("/tmp/cni.tgz", "/opt/bin/cni"), checkutil.SCPCheck)
	return cl
}
