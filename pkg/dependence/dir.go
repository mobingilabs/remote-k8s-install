package dependence

import (
	"path/filepath"

	"mobingi/ocean/pkg/constants"
	"mobingi/ocean/pkg/tools/machine"
	checkutil "mobingi/ocean/pkg/util/check"
	cmdutil "mobingi/ocean/pkg/util/cmd"
)

func GetMasterDirCommands() machine.CommandList {
	cl := machine.CommandList{}
	mkdirCheck := func(output string) bool {
		return true
	}
	mkdirCmd := cmdutil.NewMkdirAllCmd(constants.WorkDir)
	cl.Add(mkdirCmd, mkdirCheck)

	mkdirCmd = cmdutil.NewMkdirAllCmd(constants.PKIDir)
	cl.Add(mkdirCmd, mkdirCheck)

	mkdirCmd = cmdutil.NewMkdirAllCmd(filepath.Join(constants.PKIDir, "etcd"))
	cl.Add(mkdirCmd, mkdirCheck)

	return cl
}

// TODO now it is copy from getMasterDir, not true
func GetNodeDirCommands() machine.CommandList {
	cl := machine.CommandList{}
	mkdirCheck := func(output string) bool {
		return true
	}
	mkdirCmd := cmdutil.NewMkdirAllCmd(constants.WorkDir)
	cl.Add(mkdirCmd, mkdirCheck)

	mkdirCmd = cmdutil.NewMkdirAllCmd(constants.PKIDir)
	cl.Add(mkdirCmd, mkdirCheck)
	cl.Add(cmdutil.NewMkdirAllCmd("/opt/bin/cni"), checkutil.MkdirCheck)

	return cl
}
