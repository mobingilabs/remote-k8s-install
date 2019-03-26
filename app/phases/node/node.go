package node

import (
	"errors"
	"mobingi/ocean/pkg/kubernetes/service/docker"
	"path/filepath"

	"mobingi/ocean/pkg/config"
	"mobingi/ocean/pkg/constants"
	"mobingi/ocean/pkg/dependence"
	"mobingi/ocean/pkg/kubernetes/service/kubelet"
	"mobingi/ocean/pkg/log"
	"mobingi/ocean/pkg/tools/cache"
	"mobingi/ocean/pkg/tools/machine"
	checkutil "mobingi/ocean/pkg/util/check"
	cmdutil "mobingi/ocean/pkg/util/cmd"
	pkiutil "mobingi/ocean/pkg/util/pki"
)

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

	machine.AddCommandList(docker.CommandList(cfg))
	if err := machine.Run(); err != nil {
		log.Error(err)
		return err
	}
	log.Info("docker installed")

	// TODO will be move to other, and use constants not join string to path

	machine.AddCommandList(getDownloadCommands(cfg))
	if err := machine.Run(); err != nil {
		log.Error(err)
		return err
	}
	log.Info("download bin")

	machine.AddCommandList(dependence.GetNodesSetCommands(cfg))
	if err := machine.Run(); err != nil {
		log.Error(err)
		return err
	}
	log.Info("node set")

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

	machine.AddCommandList(getWriteCacert())
	if err := machine.Run(); err != nil {
		log.Error(err)
		return err
	}
	log.Info("write ca crt")

	machine.AddCommandList(kubelet.CommandList(cfg))
	if err := machine.Run(); err != nil {
		log.Error(err)
		return err
	}
	log.Info("kubelet start")

	return nil
}

func getWriteCacert() machine.CommandList {
	cl := machine.CommandList{}
	caCertByte, exists := cache.GetOne(constants.CertPrefix, pkiutil.NameForCert(constants.CACertAndKeyBaseName))
	if !exists {
		// TODO return err
		panic(errors.New("ca not get ca cert"))
	}
	writeCmd := cmdutil.NewWriteCmd(filepath.Join(constants.PKIDir, pkiutil.NameForCert(constants.CACertAndKeyBaseName)), string(caCertByte.([]byte)))
	cl.Add(writeCmd, checkutil.WriteCheck)
	return cl
}
