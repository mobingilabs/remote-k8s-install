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
	// stop kubelet
	job.AddCmd(cmd.NewSystemStopCmd(constants.KubeletService))
	job.AddCmd("docker stop `docker ps --no-trunc -aq`")
	job.AddCmd("docker rm `docker ps --no-trunc -aq`")
	job.AddCmd(cmd.NewSystemStopCmd("docker"))
	job.AddCmd("rm -rf /etc/systemd/system/kubelet.service")
	job.AddCmd("rm -rf /etc/systemd/system/kubelet.service.d")
	// delete static pod yaml file
	job.AddCmd("rm -rf /etc/kubelet.d")
	// delete kubernetes config file
	job.AddCmd("rm -rf /etc/kubernetes")
	err = machines.Run(job)
	if err != nil {
		return nil, err
	}
	return &pb.Response{Message: ""}, nil
}
