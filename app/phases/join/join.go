package join

import (
	"mobingi/ocean/pkg/config"
	"mobingi/ocean/pkg/ssh"
)

func Join(cfg *config.Config) error {
	client, err := ssh.NewClient(cfg.Nodes[0].Addr, cfg.Nodes[0].User, cfg.Nodes[0].Password)
	defer client.Close()
	if err != nil {
		return err
	}

	mkdirAll(client, cfg)

	return nil
}
