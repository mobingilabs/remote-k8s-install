package service

import (
	"mobingi/ocean/pkg/tools/machine"
)

func NewRunControlPlaneJobs(ips []string, etcdServers string) ([]*machine.Job, error) {
	apiserverJobs, err := NewRunAPIServerJobs(ips, etcdServers)
	if err != nil {
		return nil, err
	}

	controllerManagerJob, err := NewRunControllerManagerJob()
	if err != nil {
		return nil, err
	}

	schedulerJob, err := NewRunSchedulerJob()
	if err != nil {
		return nil, err
	}

	for _, v := range apiserverJobs {
		v.AddAnother(controllerManagerJob)
		v.AddAnother(schedulerJob)
	}

	return apiserverJobs, nil
}
