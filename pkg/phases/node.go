package phases

import (
	"mobingi/ocean/pkg/constants"
	"mobingi/ocean/pkg/tools/machine"
	cmdutil "mobingi/ocean/pkg/util/cmd"
	"path/filepath"
)

func NodePrepareJob(bootstrapconf []byte, certs map[string][]byte) *machine.Job {
	j := machine.NewJob("node-prepare")
	makeNodeDir(j)
	SetNodeEnv(j)
	writeNodePKI(j, certs)
	writeBootstrapconf(j, bootstrapconf)
	InstallDocker(j)

	return j
}

func makeNodeDir(j *machine.Job) {
	j.AddCmd(cmdutil.NewMkdirAllCmd(constants.WorkDir))
	j.AddCmd(cmdutil.NewMkdirAllCmd(constants.PKIDir))
	j.AddCmd(cmdutil.NewMkdirAllCmd(filepath.Join(constants.PKIDir, "etcd")))
	j.AddCmd(cmdutil.NewMkdirAllCmd("/opt/bin/cni"))
}

func SetNodeEnv(j *machine.Job) {
	j.AddCmd("swapoff -a")
	j.AddCmd(cmdutil.NewWriteCmd("/etc/sysctl.d/k8s.conf", "net.ipv4.ip_forward = 1"))
}

func writeBootstrapconf(j *machine.Job, bootstrapconf []byte) {
	j.AddCmd(cmdutil.NewWriteCmd(filepath.Join(constants.WorkDir, constants.BootstrapKubeletConfName), string(bootstrapconf)))
}

func writeNodePKI(j *machine.Job, certs map[string][]byte) {
	for k, v := range certs {
		j.AddCmd(cmdutil.NewWriteCmd(filepath.Join(constants.WorkDir, "pki", k), string(v)))
	}
}

func NodeRemoveJob() *machine.Job {
	job := machine.NewJob("delete-node")
	// stop kubelet
	job.AddCmd(cmdutil.NewSystemStopCmd(constants.KubeletService))
	// job.AddCmd("docker stop `docker ps --no-trunc -aq`")
	// job.AddCmd("docker rm `docker ps --no-trunc -aq`")
	job.AddCmd(cmdutil.NewSystemStopCmd("docker"))
	job.AddCmd("rm -rf /etc/systemd/system/kubelet.service")
	job.AddCmd("rm -rf /etc/systemd/system/kubelet.service.d")
	// delete static pod yaml file
	job.AddCmd("rm -rf /etc/kubelet.d")
	// delete kubernetes config file
	job.AddCmd("rm -rf /etc/kubernetes")
	return job
}
