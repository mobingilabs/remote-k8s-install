package app

import (
	"context"
	pb "mobingi/ocean/app/proto"
	"mobingi/ocean/pkg/config"
	preparemaster "mobingi/ocean/pkg/kubernetes/prepare/master"
	"mobingi/ocean/pkg/kubernetes/service"
	"mobingi/ocean/pkg/kubernetes/staticpod"
	"mobingi/ocean/pkg/log"
	phasesmaster "mobingi/ocean/pkg/phases/master"
	configstorage "mobingi/ocean/pkg/storage"
	"mobingi/ocean/pkg/tools/machine"
	"mobingi/ocean/pkg/util/group"
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
	configstorage.NewMongoClient()

	storage, err := configstorage.NewStorage(&configstorage.ClusterMongo{}, cfg)
	if err != nil {
		return nil, err
	}
	certs, err := storage.AllCerts()
	if err != nil {
		return nil, err
	}
	kubeconfs, err := storage.AllKubeconfs()
	if err != nil {
		return nil, err
	}

	machines := newMachines(cfg)
	log.Info("machine init")

	job := preparemaster.NewJob(cfg.DownloadBinSite, certs, kubeconfs)
	phasesmaster.InstallDocker(job)
	job.AddAnother(service.NewRunMasterKubeletJob())

	privateIPs := cfg.GetMasterPrivateIPs()
	etcdServers := service.GetEtcdServers(privateIPs)
	job.AddAnother(staticpod.NewMasterStaticPodsJob(privateIPs[0], etcdServers))

	g := group.NewGroup(len(cfg.Masters))
	for _, v := range machines {
		m := v
		g.Add(func() error {
			return m.Run(job)
		})
	}
	errs := g.Run()
	for _, v := range errs {
		if v != nil {
			return nil, err
		}
	}
	log.Info("master prepare")

	return &pb.Response{Message: ""}, nil
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
