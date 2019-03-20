package node

import (
	"errors"
	"path/filepath"

	"mobingi/ocean/pkg/config"
	"mobingi/ocean/pkg/constants"
	"mobingi/ocean/pkg/log"
	"mobingi/ocean/pkg/service/kubelet"
	"mobingi/ocean/pkg/ssh"
	"mobingi/ocean/pkg/tools/cache"
	cmdutil "mobingi/ocean/pkg/util/cmd"
)

func Start(cfg *config.Config) error {
	client, err := ssh.NewClient(cfg.Nodes[0].PublicIP, cfg.Nodes[0].User, cfg.Nodes[0].Password)
	defer client.Close()
	if err != nil {
		log.Errorf("new ssh client err:%s", err.Error())
		return err
	}
	log.Info("new ssh client sucessed")

	mkdirAll(client)
	log.Info("mkdir sucessed")

	if err := writeBootstrapConf(cfg, client); err != nil {
		log.Errorf("write bootstrap conf err:%s", err.Error())
		return err
	}
	log.Info("write bootstrap conf sucessed")

	if err := kubelet.Start(client, cfg); err != nil {
		log.Errorf("start kubelet err:%s", err.Error())
		return err
	}
	log.Info("kubelet start sucessed")

	return nil
}

func writeBootstrapConf(cfg *config.Config, c ssh.Client) error {
	// TODO make bootstarp-kubelet.conf to constants
	bootstrapConfFilename := filepath.Join(constants.WorkDir, "bootstrap-kubelet.conf")
	content, exists := cache.Get("bootstrap-kubelet.conf")
	if !exists {
		return errors.New("can not read bootstrap-kubelet.conf from cache")
	}
	c.Do(cmdutil.NewWriteCmd(bootstrapConfFilename, string(content.([]byte))))

	return nil
}
