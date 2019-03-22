package kubecontrollermanager

import (
	"mobingi/ocean/pkg/tools/machine"
	"path/filepath"

	"mobingi/ocean/pkg/config"
	"mobingi/ocean/pkg/constants"
	cmdutil "mobingi/ocean/pkg/util/cmd"
)

func CommandList(cfg *config.Config) (machine.CommandList, error) {
	serviceData, err := getServiceFile(cfg)
	if err != nil {
		return nil, err
	}

	cl := machine.CommandList{}

	writeCmd := cmdutil.NewWriteCmd(filepath.Join(constants.ServiceDir, constants.KubeControllerManagerService), string(serviceData))
	writeCheck := func(output string) bool {
		return true
	}
	cl.Add(writeCmd, writeCheck)

	startCmd := cmdutil.NewSystemStartCmd(constants.KubeControllerManagerService)
	startCheck := func(output string) bool {
		return true
	}
	cl.Add(startCmd, startCheck)

	return cl, nil
}

func getServiceFile(cfg *config.Config) ([]byte, error) {
	/*
		data, err := templateutil.Parse(serviceTemplate, newTemplateData(cfg))
		if err != nil {
			return nil, err
		}
		fmt.Println(string(data))*/

	return []byte(serviceTemplate), nil
}
