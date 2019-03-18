package join

import (
	"mobingi/ocean/pkg/config"
	"mobingi/ocean/pkg/ssh"
	cmdutil "mobingi/ocean/pkg/util/cmd"
)

func mkdirAll(c *ssh.Client, cfg *config.Config) error {
	cmd := cmdutil.NewMkdirAllCmd(cfg.WorkDir)
	c.Do(cmd)

	cmd = cmdutil.NewMkdirAllCmd(cfg.PKIDir)
	c.Do(cmd)

	// kubelet config path
	// TODO read from config
	cmd = cmdutil.NewMkdirAllCmd("/var/lib/kubelet")
	c.Do(cmd)

	return nil
}
