package master

import (
	"crypto/rsa"
	"crypto/x509"
	"fmt"
	"mobingi/ocean/pkg/kubernetes/service/etcd"
	"mobingi/ocean/pkg/kubernetes/service/kubeapiserver"
	"mobingi/ocean/pkg/kubernetes/service/kubecontrollermanager"
	"mobingi/ocean/pkg/kubernetes/service/kubescheduler"
	"mobingi/ocean/pkg/tools/cache"
	"time"

	clientset "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"

	"mobingi/ocean/pkg/certs"
	"mobingi/ocean/pkg/config"
	"mobingi/ocean/pkg/constants"
	"mobingi/ocean/pkg/dependence"
	"mobingi/ocean/pkg/kubeconfig"
	"mobingi/ocean/pkg/kubernetes/bootstrap"
	"mobingi/ocean/pkg/log"
	"mobingi/ocean/pkg/tools/machine"
	pkiutil "mobingi/ocean/pkg/util/pki"
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

	machine.AddCommandList(getDownloadCommands(cfg))
	if err := machine.Run(); err != nil {
		log.Error(err)
		return err
	}
	log.Info("master download sucess")

	machine.AddCommandList(dependence.GetMasterSetCommands())
	if err := machine.Run(); err != nil {
		log.Error(err)
		return err
	}
	log.Info("master set sucess")

	certList, err := certs.CreatePKIAssets(cfg)
	if err != nil {
		log.Errorf("create pki asstes err:%s", err)
		return err
	}
	log.Info("create pki assestes")
	machine.AddCommandList(getWriteCertsCommand(certList))
	if err := machine.Run(); err != nil {
		log.Error(err)
		return err
	}
	log.Info("write certs to disk")

	cache.Put(constants.CertPrefix, "ca.crt", certList["ca.crt"])

	etcdCommandList, err := etcd.CommandList(cfg)
	if err != nil {
		log.Error(err)
		return err
	}
	machine.AddCommandList(etcdCommandList)
	if err := machine.Run(); err != nil {
		log.Error(err)
		return err
	}
	log.Info("etcd run")

	kubeapiserverCommandList, err := kubeapiserver.CommandList(cfg)
	if err != nil {
		log.Error(err)
		return err
	}
	machine.AddCommandList(kubeapiserverCommandList)
	if err := machine.Run(); err != nil {
		log.Error(err)
		return err
	}
	log.Info("kube-apiserver run")

	kubecontrollermanagerCommandList, err := kubecontrollermanager.CommandList(cfg)
	if err != nil {
		log.Error(err)
		return err
	}
	machine.AddCommandList(kubecontrollermanagerCommandList)
	if err := machine.Run(); err != nil {
		log.Error(err)
		return err
	}
	log.Info("kube-controller-manager run")

	kubeschedulerCommandList, err := kubescheduler.CommandList(cfg)
	if err != nil {
		log.Error(err)
		return err
	}
	machine.AddCommandList(kubeschedulerCommandList)
	if err := machine.Run(); err != nil {
		log.Error(err)
		return err
	}
	log.Info("kube-scheduler run")

	caCert, caKey, err := getCaCertAndKey(certList)
	if err != nil {
		log.Error(err)
		return err
	}
	kubeconfs, err := kubeconfig.CreateKubeconf(cfg, caCert, caKey)
	if err != nil {
		log.Error(err)
		return err
	}
	log.Info("create kubeconfs")
	machine.AddCommandList(getWriteKubeconfsCommand(kubeconfs))
	if err := machine.Run(); err != nil {
		log.Error(err)
		return err
	}
	log.Info("write kubeconfs to disk")

	time.Sleep(30 * time.Second)
	k8sClient, err := newK8sClientFromConf(kubeconfs["admin.conf"])
	if err != nil {
		log.Errorf("crete k8s clinet err:%s", err.Error())
		return err
	}
	log.Info("new k8s client sucessed")

	bootstrapConf, err := bootstrap.Bootstrap(k8sClient, cfg, certList["ca.crt"])
	if err != nil {
		log.Errorf("bootstrap err:%s", err.Error())
		return err
	}
	log.Info("bootstrap done")

	cache.Put(constants.KubeconfPrefix, constants.BootstrapKubeletConfName, bootstrapConf)

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

func getCaCertAndKey(certList map[string][]byte) (*x509.Certificate, *rsa.PrivateKey, error) {
	certData, exists := certList[pkiutil.NameForCert(constants.CACertAndKeyBaseName)]
	if !exists {
		return nil, nil, fmt.Errorf("ca cert not exists in list")
	}
	cert, err := pkiutil.ParseCertPEM(certData)
	if err != nil {
		return nil, nil, err
	}

	keyData, exists := certList[pkiutil.NameForKey(constants.CACertAndKeyBaseName)]
	if !exists {
		return nil, nil, fmt.Errorf("ca key not exists in list")
	}
	key, err := pkiutil.ParsePrivateKeyPEM(keyData)
	if err != nil {
		return nil, nil, err
	}

	return cert, key, nil
}
