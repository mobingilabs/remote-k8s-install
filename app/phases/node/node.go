package node

import (
	"mobingi/ocean/pkg/kubernetes/service"
	"mobingi/ocean/pkg/kubernetes/prepare/node"
	"mobingi/ocean/pkg/log"
	"mobingi/ocean/pkg/tools/machine"
)

func Join(adminconf []byte, bootstrapconf []byte, caCert []byte, downloadBinSite string, mi *machine.MachineInfo) error {
	machine, err := machine.NewMachine(mi.PublicIP, mi.User, mi.Password)
	if err != nil {
		log.Error(err)
		return err
	}
	defer machine.Close()
	log.Info("machine init")

	// TODO load bootstrapconf from other
	err = machine.Run(node.NewJob(downloadBinSite, bootstrapconf))
	if err != nil {
		log.Panic(err)
	}
	log.Info("prepare done")

	if err := machine.Run(service.NewRunKubeletJob()); err != nil {
		log.Panic(err)
	}
	log.Info("kubelet run")

	return nil
}
