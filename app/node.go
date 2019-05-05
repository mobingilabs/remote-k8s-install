package app

import (
	"context"
	pb "mobingi/ocean/app/proto"
	"mobingi/ocean/pkg/kubernetes/service"
	"mobingi/ocean/pkg/log"
	"mobingi/ocean/pkg/phases"
	"mobingi/ocean/pkg/storage"
	"mobingi/ocean/pkg/tools/machine"
)

type node struct{}

func (n *node) Join(ctx context.Context, cfg *pb.ServerConfig) (*pb.Response, error) {
	machine, err := machine.NewMachine(cfg.PublicIP, cfg.User, cfg.Password)
	if err != nil {
		return nil, err
	}
	defer machine.Close()
	log.Info("machine init")

	storage := storage.NewStorage()
	bootstrapconf, err := storage.GetKubeconf(cfg.ClusterName, "bootstrap.conf")
	if err != nil {
		return nil, err
	}
	certs, err := storage.AllCerts(cfg.ClusterName)
	if err != nil {
		return nil, err
	}
	err = machine.Run(phases.NodePrepareJob(bootstrapconf, certs))
	if err != nil {
		return nil, err
	}
	log.Info("prepare done")

	if err := machine.Run(service.NewRunKubeletJob()); err != nil {
		return nil, err
	}
	log.Info("kubelet run")
	return &pb.Response{Message: ""}, nil
}

func (n *node) Delete(ctx context.Context, cfg *pb.ServerConfig) (*pb.Response, error) {
	machines, err := machine.NewMachine(cfg.PublicIP, cfg.User, cfg.Password)
	if err != nil {
		return nil, err
	}
	job := phases.NodeRemoveJob()
	err = machines.Run(job)
	if err != nil {
		return nil, err
	}
	return &pb.Response{Message: ""}, nil
}
