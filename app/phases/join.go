package phases

import (
	"mobingi/ocean/pkg/config"
	"mobingi/ocean/pkg/ssh"
)

func Join(cfg *config.Config) error {
	machine := cfg.GetNodeMachine()
	client, err := ssh.NewClient(machine.Addr, machine.User, machine.Password)
	defer client.Close()
	if err != nil {
		return err
	}

	return nil
}
