package app

import (
	"context"
	pb "mobingi/ocean/app/proto"
	"mobingi/ocean/pkg/constants"
	"mobingi/ocean/pkg/phases"
	"mobingi/ocean/pkg/phases/mainfests"
	"mobingi/ocean/pkg/storage"
	"mobingi/ocean/pkg/tools/machine"
	"mobingi/ocean/pkg/util/cmd"
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

	// static pod
	o := mainfests.Options{
		IP:                     cfg.PrivateIP,
		NodeName:               "node0",
		EtcdToken:              "token",
		EtcdImage:              "cnbailian/etcd:3.3.10",
		APIServerImage:         "cnbailian/kube-apiserver:v1.13.3",
		ControllerManagerImage: "cnbailian/kube-controller-manager:v1.13.3",
		SchedulerImage:         "cnbailian/kube-scheduler:v1.13.3",
		ServiceIPRange:         "10.96.0.0/12",
	}
	mainfests := mainfests.GetStaticPodMainfests(o)
	for k, v := range mainfests {
		job.AddCmd(cmd.NewWriteCmd(constants.KubeletStaticPodDir+k, string(v)))
	}

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
