package master

import (
	"path/filepath"

	"mobingi/ocean/pkg/constants"
	"mobingi/ocean/pkg/tools/machine"
	checkutil "mobingi/ocean/pkg/util/check"
	cmdutil "mobingi/ocean/pkg/util/cmd"
)

// these function may be changed the other package

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
