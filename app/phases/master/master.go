package master

import (
"sync"
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
	"mobingi/ocean/pkg/kubernetes/service"
	"mobingi/ocean/pkg/util/group"
	"mobingi/ocean/pkg/kubernetes/prepare/master"
		"mobingi/ocean/pkg/tools/kubeconf"
	"mobingi/ocean/pkg/tools/certs"
	
)

// This will be a http handler
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
	log.Info("machine init")

	g := group.NewGroup(len(cfg.Masters))
	job := master.NewJob(cfg.DownloadBinSite, certList, kubeconfs)
	for _, v := range machines {
		g.Add(func()error{
			return v.Run(job)
		})
	}
	errs := g.Run()
	for _, v := range errs {
		if v != nil {
			log.Panicf("master prepare:%s", v.Error())
		}
	}
	log.Info("master prepare")

	privateIPs := make([]string, 0, len(cfg.Masters))
	for _, v := range cfg.Masters {
		privateIPs = appned(privateIPs, v.PrivateIP)
	}

	etcdRunJobs, err := service.NewRunEtcdJobs(privateIPs)
	if err != nil {
		panic(err)
	}
	for i, v := range machines {
		g.Add(func() error {
			return v.Run(etcdRunJobs[i])
		})
	}
	errs := g.Run()
	for _, v := range errs {
		if v != nil {
			log.Panicf("etcd run:%s", v.Error())
		}
	}
	log.Info("etcd run")

	controlPlaneJobs := service.NewRunAPIServerJobs(privateIPs, service.GetEtcdServers(privateIPs))
	for i, v := range machines {
		g.Add(func() error {
			return v.Run(controlPlaneJobs[i])
		})
	}
	errs := g.Run()
	for _, v := range errs {
		if v != nil {
			log.Panicf("control plane:%s", v.Error())
		}
	}
	log.Info("control plane")

	bootstrapConf, err := bootstrap.Bootstrap(kubeconfs["admin.conf"], cfg, certList["ca.crt"])
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
