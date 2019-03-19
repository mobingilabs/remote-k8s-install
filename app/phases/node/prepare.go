package node

import (
	"mobingi/ocean/pkg/constants"
	"mobingi/ocean/pkg/ssh"
	cmdutil "mobingi/ocean/pkg/util/cmd"
)

func mkdirAll(c ssh.Client) error {
	cmd := cmdutil.NewMkdirAllCmd(constants.WorkDir)
	c.Do(cmd)

	cmd = cmdutil.NewMkdirAllCmd(constants.PKIDir)
	c.Do(cmd)

	// kubelet config path
	// TODO read from config
	cmd = cmdutil.NewMkdirAllCmd("/var/lib/kubelet")
	c.Do(cmd)

	return nil
}
