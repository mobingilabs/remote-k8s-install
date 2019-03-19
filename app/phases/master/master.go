package master

import (
	"errors"
	"fmt"
	"mobingi/ocean/pkg/constants"
	"path/filepath"

	clientset "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"

	"mobingi/ocean/pkg/certs"
	"mobingi/ocean/pkg/config"
	"mobingi/ocean/pkg/kubeconfig"
	"mobingi/ocean/pkg/kubernetes/bootstrap"
	"mobingi/ocean/pkg/log"
	"mobingi/ocean/pkg/service"
	"mobingi/ocean/pkg/ssh"
	"mobingi/ocean/pkg/tools/cache"
	cmdutil "mobingi/ocean/pkg/util/cmd"
)

func Start(cfg *config.Config) error {
	sshClient, err := ssh.NewClient(cfg.Masters[0].PublicIP, cfg.Masters[0].User, cfg.Masters[0].Password)
	if err != nil {
		log.Error(err)
		return err
	}
	log.Info("ssh client dial sucessed")
	defer sshClient.Close()

	prepare(sshClient)
	log.Info("master prepare sucessed")

	if err := certs.CreatePKIAssets(sshClient, cfg); err != nil {
		return err
	}
	log.Info("crate pki assestes sucess")

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

func mkdirAll(c ssh.Client) {
	c.Do(cmdutil.NewMkdirAllCmd(constants.WorkDir))
	c.Do(cmdutil.NewMkdirAllCmd(constants.PKIDir))
	c.Do(cmdutil.NewMkdirAllCmd(filepath.Join(constants.PKIDir, "etcd")))
}

// TODO
// docker install,set env,download binary...
func prepare(c ssh.Client) {
	mkdirAll(c)
}
