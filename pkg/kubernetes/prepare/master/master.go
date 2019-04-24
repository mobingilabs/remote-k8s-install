package kubernetes

import (
	"mobingi/ocean/pkg/constants"
	"mobingi/ocean/pkg/tools/machine"
	cmdutil "mobingi/ocean/pkg/util/cmd"
	"path/filepath"
)

func NewOneJob(kubeconfs map[string][]byte) *machine.Job {
	j := machine.NewJob("master-prepare")
	setEnv(j)
	writeKubeconfs(j, kubeconfs)
	return j
}

// NewJob will craete dir all,set machine env,write pki files to disk, write kubeconf files to disk
func NewJob(site string, certList map[string][]byte, kubeconfs map[string][]byte) *machine.Job {
	j := machine.NewJob("master-prepare")

	mkdir(j)
	setEnv(j)
	writePKI(j, certList)
	writeKubeconfs(j, kubeconfs)
	// TODO for test
	//downloadBIN(j, site)

	return j
}

func mkdir(j *machine.Job) {
	j.AddCmd(cmdutil.NewMkdirAllCmd(constants.WorkDir))
	j.AddCmd(cmdutil.NewMkdirAllCmd(constants.PKIDir))
	j.AddCmd(cmdutil.NewMkdirAllCmd(filepath.Join(constants.PKIDir, "etcd")))
}

func setEnv(j *machine.Job) {
	j.AddCmd("swapoff -a")
	// may be it need check
	j.AddCmd(cmdutil.NewWriteCmd("/etc/sysctl.d/k8s.conf", "net.ipv4.ip_forward = 1"))
}

func writePKI(j *machine.Job, certList map[string][]byte) {
	for k, v := range certList {
		j.AddCmd(cmdutil.NewWriteCmd(filepath.Join(constants.WorkDir, "pki", k), string(v)))
	}
}

func writeKubeconfs(j *machine.Job, kubeconfs map[string][]byte) {
	for k, v := range kubeconfs {
		j.AddCmd(cmdutil.NewWriteCmd(filepath.Join(constants.WorkDir, k), string(v)))
	}
}

func downloadBIN(j *machine.Job, site string) {
	j.AddCmd(cmdutil.NewCurlCmd(site, "master.tgz"))
	j.AddCmd(cmdutil.NewTarXCmd("/tmp/master.tgz", "/usr/local/bin"))
}
