package master

import (
	"mobingi/ocean/pkg/kubernetes/prepare/master"
	"sync"
	"mobingi/ocean/pkg/tools/kubeconf"
	"mobingi/ocean/pkg/tools/certs"
	"crypto/rsa"
	"crypto/x509"
	"fmt"
	"time"

	"mobingi/ocean/pkg/config"
	"mobingi/ocean/pkg/constants"
	"mobingi/ocean/pkg/dependence"
	"mobingi/ocean/pkg/tools/cache"
	"mobingi/ocean/pkg/tools/machine"
	pkiutil "mobingi/ocean/pkg/util/pki"
)

func InstallMasters(cfg *config.Config) error {
	certList, err := certs.CreatePKIAssets(cfg)
	if err != nil {
		log.Panicf("cert create:%s", err.Error())
	}
	log.Info("cert create")

	caCert, caKey, err := getCaCertAndKey(certList)
	if err != nil {
		log.Panicf("get ca cert and key :%s", err.Error())
	}
	kubeconfs, err := kubeconf.CreateKubeconf(cfg, caCert, caKey)
	if err != nil {
		log.Panicf("create kube conf :%s", err.Error())
	}
	log.Info("kubeconf create")

	machines := make([]machine.Machine, 0,len(cfg.Masters))
	for _, v := range cfg.Masters {
		machine, err := machine.NewMachine(v.PublicIP, v.User, v.Password)
		if err != nil {
			log.Panicf("new machine :%s", err.Error())
		}
	}

	errChans := make([]chan error, 0,len(cfg.Masters))
	for _, v := range machines {
		errChans = append(errChans,  v.Run(master.NewJob(cfg.DownloadBinSite, certList, kubeconfs)))
	}
	for _, v := range errChans {
	case err <- v :
		if err != nil {
			log.Panicf("master prepare:%s", err.Error())
		}
	}
	log.Info("master prepare")






	machine, err := machine.NewMachine(cfg.Masters[0].PublicIP, cfg.Masters[0].User, cfg.Masters[0].Password)
	if err != nil {
		log.Error(err)
		return err
	}
	defer machine.Close()
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

	bootstrapConf, err := bootstrap.Bootstrap(k8sClient, cfg, certList["ca.crt"])
	if err != nil {
		log.Errorf("bootstrap err:%s", err.Error())
		return err
	}
	log.Info("bootstrap done")

	cache.Put(constants.KubeconfPrefix, constants.BootstrapKubeletConfName, bootstrapConf)

	return nil
}

// it will be remove
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
