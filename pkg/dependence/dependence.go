package dependence

import (
	"mobingi/ocean/pkg/tools/machine"
)

func GetMasterSetCommand() machine.CommandList {
	cl := machine.CommandList{}
	cl.AddAnother(getSetEnvCommands())
	cl.AddAnother(getMasterBinCommands())
	return cl
}

func GetNodesSetCommand() machine.CommandList {
	cl := machine.CommandList{}
	cl.AddAnother(getSetEnvCommands())
	cl.AddAnother(getNodeBinCommands())
	cl.AddAnother(getCniBinCommands())
	return cl
}
