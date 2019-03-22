package kubeapiserver

import (
	"mobingi/ocean/pkg/tools/machine"
	"path/filepath"

	"mobingi/ocean/pkg/config"
	"mobingi/ocean/pkg/constants"
	cmdutil "mobingi/ocean/pkg/util/cmd"
	templateutil "mobingi/ocean/pkg/util/template"
)

func CommandList(cfg *config.Config) (machine.CommandList, error) {
	serviceData, err := getServiceFile(cfg)
	if err != nil {
		return nil, err
	}

	cl := machine.CommandList{}

	writeCmd := cmdutil.NewWriteCmd(filepath.Join(constants.ServiceDir, constants.KubeApiserverService), string(serviceData))
	writeCheck := func(output string) bool {
		return true
	}
	cl.Add(writeCmd, writeCheck)

	startCmd := cmdutil.NewSystemStartCmd(constants.KubeApiserverService)
	startCheck := func(output string) bool {
		return true
	}
	cl.Add(startCmd, startCheck)

	return cl, nil
}

func getServiceFile(cfg *config.Config) ([]byte, error) {
	data, err := templateutil.Parse(serviceTemplate, newTemplateData(cfg))
	if err != nil {
		return nil, err
	}

	return data, nil
}
