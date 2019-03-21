package kubelet

import (
	"path/filepath"

	"mobingi/ocean/pkg/config"
	"mobingi/ocean/pkg/constants"
	"mobingi/ocean/pkg/ssh"
	cmdutil "mobingi/ocean/pkg/util/cmd"
)

func Start(c ssh.Client, cfg *config.Config) error {
	if err := writeServiceFile(c, serviceTemplate); err != nil {
		return err
	}
	if err := writeConfigFile(c); err != nil {
		return err
	}

	if err := startService(c); err != nil {
		return err
	}

	return nil
}

func writeServiceFile(c ssh.Client, serviceData string) error {
	serviceFilepath := filepath.Join(constants.ServiceDir, constants.KubeletService)

	c.Do(cmdutil.NewWriteCmd(serviceFilepath, serviceData))
	c.Do(cmdutil.NewMkdirAllCmd(servicedDir))
	c.Do(cmdutil.NewWriteCmd(filepath.Join(servicedDir, servicedName), servicedFileContent))
	c.Do(cmdutil.NewMkdirAllCmd(configDir))
	c.Do(cmdutil.NewWriteCmd(filepath.Join(configDir, configName), configYAML))
	c.Do(cmdutil.NewWriteCmd(filepath.Join(configDir, flagsFileName), flagsContent))

	return nil
}

func writeConfigFile(c ssh.Client) error {
	cmd := cmdutil.NewMkdirAllCmd("/var/lib/kubelet")
	c.Do(cmd)

	cmd = cmdutil.NewWriteCmd("/var/lib/kubelet/config.yaml", configYAML)
	c.Do(cmd)

	return nil
}

func startService(c ssh.Client) error {
	cmd := cmdutil.NewSystemStartCmd(constants.KubeletService)
	c.DoWithoutOutput(cmd)

	return nil
}
