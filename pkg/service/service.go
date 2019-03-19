package service

import (
	"mobingi/ocean/pkg/config"
	"mobingi/ocean/pkg/service/etcd"
	"mobingi/ocean/pkg/service/kubeapiserver"
	"mobingi/ocean/pkg/service/kubecontrollermanager"
	"mobingi/ocean/pkg/service/kubelet"
	"mobingi/ocean/pkg/service/kubescheduler"
	"mobingi/ocean/pkg/ssh"
)

func Start(c ssh.Client, cfg *config.Config) error {
	if err := etcd.Start(c, cfg); err != nil {
		return err
	}

	if err := kubeapiserver.Start(c, cfg); err != nil {
		return err
	}

	if err := kubecontrollermanager.Start(c, cfg); err != nil {
		return err
	}

	if err := kubescheduler.Start(c, cfg); err != nil {
		return err
	}

	if err := kubelet.Start(c, cfg); err != nil {
		return err
	}

	return nil
}
