package kubelet

import (
	"path/filepath"

	"mobingi/ocean/pkg/ssh"
	cmdutil "mobingi/ocean/pkg/util/cmd"
)

type flagOpts struct {
	pauseImage string
}

const fileName = "/var/lib/kubelet/kubelet-flags.env"

// TODO more from config
func WriteKubeletDynamicEnvFile(mc ssh.Client) error {
	var content = `KUBELET_KUBEADM_ARGS=--cgroup-driver=systemd --network-plugin=cni --pod-infra-container-image=k8s.gcr.io/pause:3.1`
	cmd := cmdutil.NewMkdirAllCmd(filepath.Dir(fileName))
	mc.Do(cmd)
	cmd = cmdutil.NewWriteCmd(fileName, content)
	_, err := mc.Do(cmd)
	if err != nil {
		return err
	}

	return nil
}
