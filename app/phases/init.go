package phases

import (
	"errors"
	"mobingi/ocean/pkg/kubernetes/bootstrap"

	clientset "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"mobingi/ocean/pkg/tools/cache"

	"mobingi/ocean/pkg/certs"
	"mobingi/ocean/pkg/config"
	"mobingi/ocean/pkg/kubeconfig"
	"mobingi/ocean/pkg/service"
	"mobingi/ocean/pkg/ssh"
)

func Init(cfg *config.Config) error {
	machine := cfg.GetMasterMachine()
	sshClient, err := ssh.NewClient(machine.Addr, machine.User, machine.Password)
	defer sshClient.Close()
	if err != nil {
		return err
	}

	if err := certs.CreatePKIAssets(sshClient, cfg); err != nil {
		return err
	}

	if err := kubeconfig.CreateKubeconfigFiles(sshClient, cfg); err != nil {
		return err
	}

	if err := service.Start(sshClient, cfg); err != nil {
		return err
	}

	adminConf, exists := cache.Get("admin.conf")
	if !exists {
		return errors.New("no admin.conf supported")
	}

	k8sClient, err := newK8sClientFromConf(adminConf.([]byte))
	if err != nil {
		return err
	}

	err := bootstrap.Bootstrap(k8sClient, cfg)
	if err != nil {
		return err
	}

	return nil
}

func newK8sClientFromConf(conf []byte) (clientset.Interface, error) {
	config, err := clientcmd.Load(conf)
	if err != nil {
		return nil, err
	}

	clientConfig, err := clientcmd.NewDefaultClientConfig(*config, &clientcmd.ConfigOverrides{}).ClientConfig()
	if err != nil {
		return nil, err
	}

	client, err := clientset.NewForConfig(clientConfig)
	if err != nil {
		return nil, err
	}

	return client, nil
}
