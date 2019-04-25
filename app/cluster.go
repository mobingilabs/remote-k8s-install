package app

import (
	"context"
	pb "mobingi/ocean/app/proto"
	"mobingi/ocean/pkg/config"
	"mobingi/ocean/pkg/storage"
)

type cluster struct{}

func (c *cluster) Init(ctx context.Context, ccfg *pb.ClusterConfig) (*pb.Response, error) {
	var masters []config.Machine
	for _, v := range ccfg.Masters {
		machine := config.Machine{
			PublicIP:  v.PublicIP,
			PrivateIP: v.PrivateIP,
			User:      v.User,
			Password:  v.Password,
		}
		masters = append(masters, machine)
	}
	cfg := &config.Config{
		ClusterName:      ccfg.ClusterName,
		AdvertiseAddress: ccfg.AdvertiseAddress,
		DownloadBinSite:  ccfg.DownloadBinSite,
		PublicIP:         ccfg.PublicIP,
		Masters:          []config.Machine(masters),
	}

	// TODO Move to main init func
	storage.NewMongoClient()

	_, err := storage.NewStorage(&storage.ClusterMongo{}, cfg)
	if err != nil {
		return nil, err
	}

	// sans := cfg.GetSANs()
	// _, err := certs.CreatePKIAssets(cfg.AdvertiseAddress, cfg.PublicIP, sans)

	// db, err := sql.Open("mysql", "root:123456789@/kubeconf")
	// if err != nil {
	// 	log.Panicf("conn: %s", err.Error())
	// }
	// defer db.Close()

	// certList, err := getConfigBySql(db, "certs", func() (map[string][]byte, error) {
	// 	sans := cfg.GetSANs()
	// 	return certs.CreatePKIAssets(cfg.AdvertiseAddress, cfg.PublicIP, sans)
	// })
	// if err != nil {
	// 	log.Panicf("cert create:%s", err.Error())
	// }

	// caCert, caKey, err := getCaCertAndKey(certList)
	// if err != nil {
	// 	log.Panicf("get ca cert and key :%s", err.Error())
	// }
	// kubeconfs, err := kubeconf.CreateKubeconf(cfg, caCert, caKey)

	// kubeconfs, err := getConfigBySql(db, "kubeconfs", func() (map[string][]byte, error) {
	// 	caCert, caKey, err := getCaCertAndKey(certList)
	// 	if err != nil {
	// 		log.Panicf("get ca cert and key :%s", err.Error())
	// 	}
	// 	return kubeconf.CreateKubeconf(cfg, caCert, caKey)
	// })
	// if err != nil {
	// 	log.Panicf("cert create:%s", err.Error())
	// }

	// log.Info("kubeconf create")

	// machines := newMachines(cfg)
	// privateIPs := cfg.GetMasterPrivateIPs()
	// runEtcdCluster(machines, privateIPs)
	// EtcdServers = service.GetEtcdServers(privateIPs)

	// log.Info("Etcd cluster create")

	return &pb.Response{Message: ""}, nil
}
