package app

import (
	"mobingi/ocean/app/phases/master"
	"mobingi/ocean/app/phases/node"
	"mobingi/ocean/pkg/config"
	"mobingi/ocean/pkg/constants"
	"mobingi/ocean/pkg/kubernetes/bootstrap"
	"mobingi/ocean/pkg/log"
	"mobingi/ocean/pkg/tools/cache"
	"mobingi/ocean/pkg/tools/machine"
)

func Start() error {
	cfg, err := config.LoadConfigFromFile("config.yaml")
	if err != nil {
		return err
	}
	if err := master.InstallMasters(cfg); err != nil {
		return err
	}

	adminConf, _ := cache.GetOne(constants.KubeconfPrefix, "admin.conf")

	bootstrapconf, err := bootstrap.Bootstrap(adminConf.([]byte))
	if err != nil {
		log.Panic(err)
	}

	mi := &machine.MachineInfo{
		PublicIP: cfg.Nodes[0].PublicIP,
		User:     cfg.Nodes[0].User,
		Password: cfg.Nodes[0].Password,
	}

	if err := node.Join(adminConf.([]byte), bootstrapconf, cfg.DownloadBinSite, mi); err != nil {
		return err
	}
	mi.PublicIP = cfg.Nodes[1].PublicIP
	mi.User = cfg.Nodes[1].User
	mi.Password = cfg.Nodes[0].Password
	if err := node.Join(adminConf.([]byte), bootstrapconf, cfg.DownloadBinSite, mi); err != nil {
		return err
	}

	return nil
}
