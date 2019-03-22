package node

import (
	"errors"

	"mobingi/ocean/pkg/config"
	"mobingi/ocean/pkg/constants"
	"mobingi/ocean/pkg/dependence"
	"mobingi/ocean/pkg/kubernetes/service/kubelet"
	"mobingi/ocean/pkg/log"
	"mobingi/ocean/pkg/tools/cache"
	"mobingi/ocean/pkg/tools/machine")

func Start(cfg *config.Config) error {
	machine, err := machine.NewMachine(cfg.Nodes[0].PublicIP, cfg.Nodes[0].User, cfg.Nodes[0].Password)
	if err != nil {
		log.Error(err)
		return err
	}
	defer machine.DisConnect()
	log.Info("machine init")

	machine.AddCommandList(dependence.GetNodeDirCommands())
	if err := machine.Run(); err != nil {
		log.Error(err)
		return err
	}
	log.Info("node create dirs")

	bootstrapConfByte, exists := cache.GetOne(constants.KubeconfPrefix, constants.BootstrapKubeletConfName)
	if !exists {
		return errors.New("bootstarp conf not in cache")
	}
	machine.AddCommandList(getWriteBootstrapCommands(bootstrapConfByte.([]byte)))
	if err := machine.Run(); err != nil {
		log.Error(err)
		return err
	}
	log.Info("write bootstrap conf to disk")

	machine.AddCommandList(kubelet.CommandList(cfg))
	if err := machine.Run(); err != nil {
		log.Error(err)
		return err
	}
	log.Info("kubelet start")

	return nil
}
