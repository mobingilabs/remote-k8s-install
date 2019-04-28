package app

import (
	"context"
	pb "mobingi/ocean/app/proto"
	"mobingi/ocean/pkg/constants"
	"mobingi/ocean/pkg/tools/machine"
	"mobingi/ocean/pkg/util/cmd"
)

type master struct{}

func (m *master) Join(ctx context.Context, cfg *pb.ServerConfig) (*pb.Response, error) {
	return nil, nil
	// job := preparemaster.NewOneJob(phasesmaster.Kubeconfs)
	// machine, err := machine.NewMachine(cfg.PublicIP, cfg.User, cfg.Password)
	// if err != nil {
	// 	return nil, err
	// }

	// apiserverJob, err := service.NewOneRunAPIServerJob(cfg.PrivateIP, phasesmaster.EtcdServers, phasesmaster.MasterCommonConfig.AdvertiseAddress)
	// if err != nil {
	// 	return nil, err
	// }
	// controllerManagerJob, err := service.NewRunControllerManagerJob()
	// if err != nil {
	// 	return nil, err
	// }
	// schedulerJob, err := service.NewRunSchedulerJob()
	// if err != nil {
	// 	return nil, err
	// }
	// job.AddAnother(apiserverJob)
	// job.AddAnother(controllerManagerJob)
	// job.AddAnother(schedulerJob)

	// err = machine.Run(job)
	// if err != nil {
	// 	fmt.Print(err.Error())
	// 	return nil, err
	// }

	// return &pb.Response{Message: ""}, nil
}

func (m *master) Delete(ctx context.Context, cfg *pb.ServerConfig) (*pb.Response, error) {
	machines, err := machine.NewMachine(cfg.PublicIP, cfg.User, cfg.Password)
	if err != nil {
		return nil, err
	}

	job := machine.NewJob("delete-master")
	job.AddCmd(cmd.NewSystemStopCmd(constants.KubeApiserverService))
	job.AddCmd(cmd.NewSystemStopCmd(constants.KubeControllerManagerService))
	job.AddCmd(cmd.NewSystemStopCmd(constants.KubeSchedulerService))
	job.AddCmd("rm /etc/systemd/system/kube-apiserver.service")
	job.AddCmd("rm /etc/systemd/system/kube-controller-manager.service")
	job.AddCmd("rm /etc/systemd/system/kube-scheduler.service")
	job.AddCmd("rm -rf /etc/kubernetes")
	err = machines.Run(job)
	if err != nil {
		return nil, err
	}
	return &pb.Response{Message: ""}, nil
}
