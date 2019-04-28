package staticpod

import (
	"mobingi/ocean/pkg/constants"
	"mobingi/ocean/pkg/log"
	"mobingi/ocean/pkg/tools/machine"
	cmdutil "mobingi/ocean/pkg/util/cmd"
	"path/filepath"
)

func NewMasterStaticPodsJob(ip, etcdServers string) *machine.Job {
	j := machine.NewJob("master-static-pod")

	etcdYaml, err := getEtcdStaticPodFile(ip)
	if err != nil {
		log.Panic(err.Error())
	}
	j.AddCmd(cmdutil.NewWriteCmd(filepath.Join(constants.KubeletStaticPodDir, constants.EtcdPod), string(etcdYaml)))
	apiserverYaml, err := getAPIServerStaticPodFile(ip, etcdServers)
	if err != nil {
		log.Panic(err.Error())
	}
	j.AddCmd(cmdutil.NewWriteCmd(filepath.Join(constants.KubeletStaticPodDir, constants.KubeApiserverPod), string(apiserverYaml)))
	j.AddCmd(cmdutil.NewWriteCmd(filepath.Join(constants.KubeletStaticPodDir, constants.KubeControllerManagerPod), controllerManagerYaml))
	j.AddCmd(cmdutil.NewWriteCmd(filepath.Join(constants.KubeletStaticPodDir, constants.KubeSchedulerPod), schedulerYaml))

	return j
}
