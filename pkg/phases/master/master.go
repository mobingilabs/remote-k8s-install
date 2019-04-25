package master

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"fmt"

	"database/sql"
	"mobingi/ocean/pkg/config"
	"mobingi/ocean/pkg/constants"
	"mobingi/ocean/pkg/kubernetes/service"
	"mobingi/ocean/pkg/log"
	"mobingi/ocean/pkg/tools/certs"
	"mobingi/ocean/pkg/tools/kubeconf"
	"mobingi/ocean/pkg/tools/machine"
	"mobingi/ocean/pkg/util/group"
	pkiutil "mobingi/ocean/pkg/util/pki"

	_ "github.com/go-sql-driver/mysql"
)

type configSql struct {
	id      int
	name    string
	context string
}

var CertList map[string][]byte
var Kubeconfs map[string][]byte
var MasterCommonConfig *config.Config
var EtcdServers string

func Init(cfg *config.Config) error {
	MasterCommonConfig = cfg
	db, err := sql.Open("mysql", "root:123456789@/kubeconf")
	if err != nil {
		log.Panicf("conn: %s", err.Error())
	}
	defer db.Close()

	CertList, err = getConfigBySql(db, "certs", func() (map[string][]byte, error) {
		sans := cfg.GetSANs()
		return certs.CreatePKIAssets(cfg.AdvertiseAddress, cfg.PublicIP, sans)
	})
	if err != nil {
		log.Panicf("cert create:%s", err.Error())
	}

	Kubeconfs, err = getConfigBySql(db, "kubeconfs", func() (map[string][]byte, error) {
		caCert, caKey, err := getCaCertAndKey(CertList)
		if err != nil {
			log.Panicf("get ca cert and key :%s", err.Error())
		}
		return kubeconf.CreateKubeconf(cfg, caCert, caKey)
	})
	if err != nil {
		log.Panicf("cert create:%s", err.Error())
	}

	log.Info("kubeconf create")

	machines := newMachines(cfg)
	privateIPs := cfg.GetMasterPrivateIPs()
	runEtcdCluster(machines, privateIPs)
	EtcdServers = service.GetEtcdServers(privateIPs)

	log.Info("Etcd cluster create")

	return nil
}

func emptyDBConfig(db *sql.DB) {
	_, err := db.Exec("DELETE FROM config")
	if err != nil {
		log.Panic(err.Error())
	}
}

func getConfigBySql(db *sql.DB, name string, callback func() (map[string][]byte, error)) (map[string][]byte, error) {
	var config map[string][]byte
	rows, err := db.Query("SELECT * FROM config WHERE name = ? LIMIT 1", name)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var context configSql
		if err := rows.Scan(&context.id, &context.name, &context.context); err != nil {
			log.Panic(err.Error())
		}
		err := json.Unmarshal([]byte(context.context), &config)
		if err != nil {
			log.Panic(err.Error())
		}
	}
	if len(config) == 0 {
		config, err = callback()
		if err != nil {
			return nil, err
		}
		log.Info("cert create")

		configJSON, err := json.Marshal(config)
		if err != nil {
			return nil, err
		}
		_, err = db.Exec(
			"INSERT INTO config (name, context) VALUES (?, ?)",
			name,
			string(configJSON),
		)
		if err != nil {
			return nil, err
		}
	}
	return config, nil
}

// This will be a http handler
// func InstallMasters(cfg *config.Config) error {
// 	sans := cfg.GetSANs()
// 	certList, err := certs.CreatePKIAssets(cfg.AdvertiseAddress, cfg.PublicIP, sans)
// 	if err != nil {
// 		log.Panicf("cert create:%s", err.Error())
// 	}
// 	log.Info("cert create")

// 	caCert, caKey, err := getCaCertAndKey(certList)
// 	if err != nil {
// 		log.Panicf("get ca cert and key :%s", err.Error())
// 	}
// 	kubeconfs, err := kubeconf.CreateKubeconf(cfg, caCert, caKey)
// 	if err != nil {
// 		log.Panicf("create kube conf :%s", err.Error())
// 	}
// 	log.Info("kubeconf create")
// 	// TODO we will put confs to store, not cache
// 	cache.Put(constants.KubeconfPrefix, "admin.conf", kubeconfs["admin.conf"])

// 	machines := newMachines(cfg)
// 	log.Info("machine init")

// 	g := group.NewGroup(len(cfg.Masters))
// 	job := preparemaster.NewJob(cfg.DownloadBinSite, certList, kubeconfs)
// 	for _, v := range machines {
// 		m := v
// 		g.Add(func() error {
// 			return m.Run(job)
// 		})
// 	}
// 	errs := g.Run()
// 	for _, v := range errs {
// 		if v != nil {
// 			log.Panicf("master prepare:%s", v.Error())
// 		}
// 	}
// 	log.Info("master prepare")

// 	privateIPs := cfg.GetMasterPrivateIPs()
// 	runEtcdCluster(machines, privateIPs)

// 	etcdServers := service.GetEtcdServers(privateIPs)
// 	runControlPlane(machines, privateIPs, etcdServers, cfg.AdvertiseAddress)

// 	// TODO wait for services up
// 	time.Sleep(time.Second)

// 	return nil
// }

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

func newMachines(cfg *config.Config) []machine.Machine {
	machines := make([]machine.Machine, 0, len(cfg.Masters))
	for _, v := range cfg.Masters {
		machine, err := machine.NewMachine(v.PublicIP, v.User, v.Password)
		if err != nil {
			log.Panicf("new machine :%s", err.Error())
		}
		machines = append(machines, machine)
	}

	return machines
}

func runEtcdCluster(machines []machine.Machine, privateIPs []string) {
	etcdRunJobs, err := service.NewRunEtcdJobs(privateIPs, CertList)
	if err != nil {
		panic(err)
	}

	g := group.NewGroup(len(machines))
	for i, v := range machines {
		m := v
		j := i
		g.Add(func() error {
			return m.Run(etcdRunJobs[j])
		})
	}
	// TODO we will design a error list type for check easily
	errs := g.Run()
	for _, v := range errs {
		if v != nil {
			log.Panicf("etcd run:%s", v.Error())
		}
	}
	log.Info("etcd run")
}

func runControlPlane(machines []machine.Machine, privateIPs []string, etcdServers, advertiseAddress string) {
	controlPlaneJobs, err := service.NewRunControlPlaneJobs(privateIPs, etcdServers, advertiseAddress)
	if err != nil {
		log.Panic(err)
	}

	g := group.NewGroup(len(machines))
	for i, v := range machines {
		m := v
		j := i
		g.Add(func() error {
			return m.Run(controlPlaneJobs[j])
		})
	}
	errs := g.Run()
	for _, v := range errs {
		if v != nil {
			log.Panicf("control plane:%s", v.Error())
		}
	}
	log.Info("control plane")
}
