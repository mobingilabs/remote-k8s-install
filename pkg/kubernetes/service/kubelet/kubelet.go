package kubelet

import (
	"mobingi/ocean/pkg/tools/machine"
	checkutil "mobingi/ocean/pkg/util/check"
	"path/filepath"

	"mobingi/ocean/pkg/config"
	"mobingi/ocean/pkg/constants"
	cmdutil "mobingi/ocean/pkg/util/cmd"
)

func CommandList(cfg *config.Config) machine.CommandList {
	cl := machine.CommandList{}

	mkdirCmd0 := cmdutil.NewMkdirAllCmd(filepath.Join(constants.ServiceDir, servicedDir))
	cl.Add(mkdirCmd0, checkutil.MkdirCheck)

	mkdirCmd1 := cmdutil.NewMkdirAllCmd(configDir)
	cl.Add(mkdirCmd1, checkutil.MkdirCheck)

	writeCmd1 := cmdutil.NewWriteCmd(filepath.Join(constants.ServiceDir, constants.KubeletService), serviceTemplate)
	writeCheck1 := func(output string) bool {
		return true
	}
	cl.Add(writeCmd1, writeCheck1)

	writeCmd2 := cmdutil.NewWriteCmd(filepath.Join(constants.ServiceDir, servicedDir, servicedName), servicedFileContent)
	writeCheck2 := func(output string) bool {
		return true
	}
	cl.Add(writeCmd2, writeCheck2)

	writeCmd3 := cmdutil.NewWriteCmd(filepath.Join(configDir, configName), configYAML)
	writeCheck3 := func(output string) bool {
		return true
	}
	cl.Add(writeCmd3, writeCheck3)

	writeCmd4 := cmdutil.NewWriteCmd(filepath.Join(configDir, flagsFileName), flagsContent)
	writeCheck4 := func(output string) bool {
		return true
	}
	cl.Add(writeCmd4, writeCheck4)

	startCmd := cmdutil.NewSystemStartCmd(constants.KubeApiserverService)
	startCheck := func(output string) bool {
		return true
	}
	cl.Add(startCmd, startCheck)

	return cl
}
