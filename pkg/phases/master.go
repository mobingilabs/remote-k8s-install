package phases

import (
	"mobingi/ocean/pkg/constants"
	"mobingi/ocean/pkg/tools/machine"
	cmdutil "mobingi/ocean/pkg/util/cmd"
	"path/filepath"
)

func MasterPrepareJob(certs map[string][]byte, kubeconfs map[string][]byte) *machine.Job {
	j := machine.NewJob("master-prepare")
	writeMasterPKI(j, certs)
	writeMasterKubeconfs(j, kubeconfs)
	j.AddCmd(cmdutil.NewSystemStartCmd("docker"))
	j.AddCmd(cmdutil.NewSystemStartCmd("kubelet"))
	return j
}

func writeMasterPKI(j *machine.Job, certs map[string][]byte) {
	for k, v := range certs {
		j.AddCmd(cmdutil.NewWriteCmd(filepath.Join(constants.WorkDir, "pki", k), string(v)))
	}
}

func writeMasterKubeconfs(j *machine.Job, kubeconfs map[string][]byte) {
	for k, v := range kubeconfs {
		j.AddCmd(cmdutil.NewWriteCmd(filepath.Join(constants.WorkDir, k), string(v)))
	}
}

func MasterRemoveJob() *machine.Job {
	job := machine.NewJob("delete-master")
	// stop kubelet
	job.AddCmd(cmdutil.NewSystemStopCmd(constants.KubeletService))
	job.AddCmd("docker stop `docker ps --no-trunc -aq`")
	job.AddCmd("docker rm `docker ps --no-trunc -aq`")
	job.AddCmd(cmdutil.NewSystemStopCmd("docker"))
	// job.AddCmd("rm -rf /etc/systemd/system/kubelet.service")
	// job.AddCmd("rm -rf /etc/systemd/system/kubelet.service.d")
	// delete static pod yaml file
	// job.AddCmd("rm -rf /etc/kubelet.d")
	// delete kubernetes config file
	// job.AddCmd("rm -rf /etc/kubernetes")
	return job
}
