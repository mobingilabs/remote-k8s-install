package node

import (
	"errors"
	"path/filepath"

	"mobingi/ocean/pkg/config"
	"mobingi/ocean/pkg/service/kubelet"
	"mobingi/ocean/pkg/ssh"
	"mobingi/ocean/pkg/tools/cache"
	cmdutil "mobingi/ocean/pkg/util/cmd"
)

func Start(cfg *config.Config) error {
	client, err := ssh.NewClient(cfg.Nodes[0].Addr, cfg.Nodes[0].User, cfg.Nodes[0].Password)
	defer client.Close()
	if err != nil {
		return err
	}

	mkdirAll(client, cfg)

	if err := writeBootstrapConf(cfg, client); err != nil {
		return err
	}

	if err := kubelet.Start(client, cfg); err != nil {
		return err
	}

	return nil

}

func writeBootstrapConf(cfg *config.Config, c *ssh.Client) error {
	// TODO make bootstarp-kubelet.conf to constants
	bootstrapConfFilename := filepath.Join(cfg.WorkDir, "bootstrap-kubelet.conf")
	content, exists := cache.Get("bootstrap-kubelet.conf")
	if !exists {
		return errors.New("can not read bootstrap-kubelet.conf from cache")
	}
	cmdutil.NewWriteCmd(bootstrapConfFilename, content.(string))

	return nil
}
