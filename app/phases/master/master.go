package master

import (
	"mobingi/ocean/pkg/dependence"
	"errors"
	"path/filepath"

	clientset "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"

	"mobingi/ocean/pkg/certs"
	"mobingi/ocean/pkg/config"
	"mobingi/ocean/pkg/constants"
	"mobingi/ocean/pkg/kubeconfig"
	"mobingi/ocean/pkg/kubernetes/bootstrap"
	"mobingi/ocean/pkg/log"
	"mobingi/ocean/pkg/service"
	"mobingi/ocean/pkg/ssh"
	"mobingi/ocean/pkg/tools/cache"
	"mobingi/ocean/pkg/tools/machine"
	cmdutil "mobingi/ocean/pkg/util/cmd"
)

func Start(cfg *config.Config) error {
	machine, err := machine.NewMachine(cfg.Masters[0].PublicIP, cfg.Masters[0].User, cfg.Masters[0].Password)
	if err != nil {
		log.Error(err)
		return err
	}
	defer machine.DisConnect()
	log.Info("machine init")

	machine.AddCommandList(dependence.GetMasterDirCommands())
	if err := machine.Run(); err != nil {
		log.Error(err)
		return err
	}
	log.Info("master create dirs")
	
	if err := certs.CreatePKIAssets(sshClient, cfg); err != nil {
		return err
	}
	log.Info("crate pki assestes sucessed")

	if err := kubeconfig.CreateKubeconfigFiles(sshClient, cfg); err != nil {
		return err
	}
	log.Info("create kube config files sucessed")

	if err := service.Start(sshClient, cfg); err != nil {
		return err
	}
	log.Info("service started")

	adminConf, exists := cache.Get("admin.conf")
	if !exists {
		log.Error("no admin.conf supported from cache")
		return errors.New("no admin.conf supported")
	}

	k8sClient, err := newK8sClientFromConf(adminConf.([]byte))
	if err != nil {
		log.Errorf("crete k8s clinet err:%s", err.Error())
		return err
	}
	log.Info("new k8s client sucessed")

	err = bootstrap.Bootstrap(k8sClient, cfg)
	if err != nil {
		log.Errorf("bootstrap err:%s", err.Error())
		return err
	}
	log.Info("bootstrap suecssed")

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

// TODO list to init
func mkdirAll(c ssh.Client) {
	c.Do(cmdutil.NewMkdirAllCmd(constants.WorkDir))
	c.Do(cmdutil.NewMkdirAllCmd(constants.PKIDir))
	c.Do(cmdutil.NewMkdirAllCmd(filepath.Join(constants.PKIDir, "etcd")))
	c.Do(cmdutil.NewMkdirAllCmd(constants.ETCDDataDir))
}

// TODO
func prepare(c ssh.Client) {
	mkdirAll(c)
}
