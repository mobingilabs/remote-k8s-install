package etcd

import (
	"path/filepath"

	"mobingi/ocean/pkg/config"
	"mobingi/ocean/pkg/constants"
	"mobingi/ocean/pkg/tools/machine"
	cmdutil "mobingi/ocean/pkg/util/cmd"
	templateutil "mobingi/ocean/pkg/util/template"
)

func CommandList(cfg *config.Config) (machine.CommandList, error) {
	cl := machine.CommandList{}
	serviceData, err := getServiceFile(cfg)
	if err != nil {
		return nil, err
	}

	writeCmd := cmdutil.NewWriteCmd(filepath.Join(constants.ServiceDir, constants.EtcdService), string(serviceData))
	writeCheck := func(output string) bool {
		return true
	}
	cl.Add(writeCmd, writeCheck)

	startCmd := cmdutil.NewSystemStartCmd(constants.EtcdService)
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
