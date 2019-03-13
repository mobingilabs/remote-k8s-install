package phases

import (
	"mobingi/ocean/pkg/certs"
	"mobingi/ocean/pkg/config"
	"mobingi/ocean/pkg/kubeconfig"
	"mobingi/ocean/pkg/service"
	"mobingi/ocean/pkg/ssh"
)

func Init(cfg *config.Config) error {
	machine := getMasterMachine(cfg)
	client, err := ssh.NewClient(machine.Addr, machine.User, machine.Password)
	defer client.Close()
	if err != nil {
		return err
	}

	if err := certs.CreatePKIAssets(client, cfg); err != nil {
		return err
	}

	if err := kubeconfig.CreateKubeconfigFiles(client, cfg); err != nil {
		return err
	}

	if err := service.Start(client, cfg); err != nil {
		return err
	}

	return nil
}

func getMasterMachine(cfg *config.Config) *config.Machine {
	for _, v := range cfg.Machines {
		if v.Role == config.MasterRole {
			return &v
		}
	}

	return nil
}
