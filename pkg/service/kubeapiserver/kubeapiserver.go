package kubeapiserver

import (
	"path/filepath"

	"mobingi/ocean/pkg/config"
	"mobingi/ocean/pkg/constants"
	"mobingi/ocean/pkg/ssh"
	cmdutil "mobingi/ocean/pkg/util/cmd"
	templateutil "mobingi/ocean/pkg/util/template"
)

func Start(c *ssh.Client, cfg *config.Config) error {
	serviceData, err := getServiceFile(cfg)
	if err != nil {
		return err
	}

	if err := writeServiceFile(c, string(serviceData)); err != nil {
		return err
	}

	if err := startService(c); err != nil {
		return err
	}

	return nil
}

func getServiceFile(cfg *config.Config) ([]byte, error) {
	data, err := templateutil.Parse(serviceTemplate, newTemplateData(cfg))
	if err != nil {
		return nil, err
	}

	return data, nil
}

func writeServiceFile(c *ssh.Client, serviceData string) error {
	serviceFilepath := filepath.Join(constants.ServiceDir, constants.KubeApiserverService)

	cmd := cmdutil.NewWriteCmd(serviceFilepath, serviceData)
	// TODO check output exec result, ok or false
	c.Do(cmd)

	return nil
}

func startService(c *ssh.Client) error {
	cmd := cmdutil.NewSystemStartCmd(constants.KubeApiserverService)
	c.Do(cmd)

	return nil
}
