package app

import (
	"context"
	pb "mobingi/ocean/app/proto"
	"mobingi/ocean/pkg/kubernetes/staticpod"
	"mobingi/ocean/pkg/phases"
	"mobingi/ocean/pkg/storage"
	"mobingi/ocean/pkg/tools/machine"
)

type master struct{}

func (m *master) Join(ctx context.Context, cfg *pb.ServerConfig) (*pb.Response, error) {
	storage := storage.NewStorage()
	certs, err := storage.AllCerts(cfg.ClusterName)
	if err != nil {
		return nil, err
	}
	kubeconfs, err := storage.AllKubeconfs(cfg.ClusterName)
	if err != nil {
		return nil, err
	}

	machine, err := machine.NewMachine(cfg.PublicIP, cfg.User, cfg.Password)
	if err != nil {
		return nil, err
	}
	job := phases.MasterPrepareJob(certs, kubeconfs)
	etcdServers, err := storage.GetEtcdServers(cfg.ClusterName)
	job.AddAnother(staticpod.NewMasterStaticPodsJob(cfg.PrivateIP, etcdServers))
	err = machine.Run(job)
	if err != nil {
		return nil, err
	}

	return &pb.Response{Message: ""}, nil
}

func (m *master) Delete(ctx context.Context, cfg *pb.ServerConfig) (*pb.Response, error) {
	machines, err := machine.NewMachine(cfg.PublicIP, cfg.User, cfg.Password)
	if err != nil {
		return nil, err
	}
	job := phases.MasterRemoveJob()
	err = machines.Run(job)
	if err != nil {
		return nil, err
	}
	return &pb.Response{Message: ""}, nil
}
