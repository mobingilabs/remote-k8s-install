package node

import (
	"mobingi/ocean/pkg/constants"
	"mobingi/ocean/pkg/ssh"
	cmdutil "mobingi/ocean/pkg/util/cmd"
)

func mkdirAll(c ssh.Client) error {
	c.Do(cmdutil.NewMkdirAllCmd(constants.WorkDir))
	c.Do(cmdutil.NewMkdirAllCmd(constants.PKIDir))
	// kubelet config path
	// TODO read from config
	c.Do(cmdutil.NewMkdirAllCmd("/var/lib/kubelet"))

	return nil
}
