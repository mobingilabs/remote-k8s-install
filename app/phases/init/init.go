package init

import (
	"errors"
	"fmt"
	"mobingi/ocean/pkg/kubernetes/bootstrap"

	"mobingi/ocean/pkg/tools/cache"

	clientset "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"

	"mobingi/ocean/pkg/certs"
	"mobingi/ocean/pkg/config"
	"mobingi/ocean/pkg/kubeconfig"
	"mobingi/ocean/pkg/service"
	"mobingi/ocean/pkg/ssh"
)

func Init(cfg *config.Config) error {
	sshClient, err := ssh.NewClient(cfg.Masters[0].Addr, cfg.Masters[1].User, cfg.Masters[2].Password)
	if err != nil {
		return err
	}
	defer sshClient.Close()

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
	fmt.Println(adminConf.(string))
	if !exists {
		return errors.New("no admin.conf supported")
	}

	k8sClient, err := newK8sClientFromConf(adminConf.([]byte))
	if err != nil {
		return err
	}

	err = bootstrap.Bootstrap(k8sClient, cfg)
	if err != nil {
		fmt.Println(err)
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
