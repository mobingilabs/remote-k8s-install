package master

import (
	"crypto/rsa"
	"crypto/x509"
	"fmt"

	"mobingi/ocean/pkg/config"
	"mobingi/ocean/pkg/constants"
	preparemaster "mobingi/ocean/pkg/kubernetes/prepare/master"
	"mobingi/ocean/pkg/kubernetes/service"
	"mobingi/ocean/pkg/log"
	"mobingi/ocean/pkg/tools/cache"
	"mobingi/ocean/pkg/tools/certs"
	"mobingi/ocean/pkg/tools/kubeconf"
	"mobingi/ocean/pkg/tools/machine"
	"mobingi/ocean/pkg/util/group"
	pkiutil "mobingi/ocean/pkg/util/pki"
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
	// TODO we will put confs to store, not cache
	cache.Put(constants.KubeconfPrefix, "admin.conf", kubeconfs["admin.conf"])

	machines := make([]machine.Machine, 0, len(cfg.Masters))
	for _, v := range cfg.Masters {
		machine, err := machine.NewMachine(v.PublicIP, v.User, v.Password)
		if err != nil {
			log.Panicf("new machine :%s", err.Error())
		}
		machines = append(machines, machine)
	}
	log.Info("machine init")

	g := group.NewGroup(len(cfg.Masters))
	job := preparemaster.NewJob(cfg.DownloadBinSite, certList, kubeconfs)
	for _, v := range machines {
		g.Add(func() error {
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
		privateIPs = append(privateIPs, v.PrivateIP)
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
	errs = g.Run()
	for _, v := range errs {
		if v != nil {
			log.Panicf("etcd run:%s", v.Error())
		}
	}
	log.Info("etcd run")

	controlPlaneJobs, err := service.NewRunControlPlaneJobs(privateIPs, service.GetEtcdServers(privateIPs))
	if err != nil {
		log.Panic(err)
	}
	for i, v := range machines {
		g.Add(func() error {
			return v.Run(controlPlaneJobs[i])
		})
	}
	errs = g.Run()
	for _, v := range errs {
		if v != nil {
			log.Panicf("control plane:%s", v.Error())
		}
	}
	log.Info("control plane")

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
