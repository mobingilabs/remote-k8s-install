package node

import (
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
