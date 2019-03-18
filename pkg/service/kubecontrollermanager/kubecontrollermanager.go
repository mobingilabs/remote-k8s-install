package kubecontrollermanager

import (
	"path/filepath"

	"mobingi/ocean/pkg/config"
	"mobingi/ocean/pkg/constants"
	"mobingi/ocean/pkg/ssh"
	cmdutil "mobingi/ocean/pkg/util/cmd"
)

func Start(c *ssh.Client, cfg *config.Config) error {
	serviceData, err := getServiceFile(cfg)
	if err != nil {
		return err
	}
	if err := writeServiceFile(c, string(serviceData)); err != nil {
		return err
	}

	return startService(c)
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

func writeServiceFile(c *ssh.Client, serviceData string) error {
	serviceFilepath := filepath.Join(constants.ServiceDir, constants.KubeControllerManagerService)

	cmd := cmdutil.NewWriteCmd(serviceFilepath, serviceData)
	// TODO check output exec result, ok or false
	_, err := c.Do(cmd)
	if err != nil {
		return err
	}

	return nil
}

func startService(c *ssh.Client) error {
	cmd := cmdutil.NewSystemStartCmd(constants.KubeControllerManagerService)
	c.DoWithoutOutput(cmd)

	return nil
}
