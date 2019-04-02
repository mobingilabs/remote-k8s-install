package service

import (
	"mobingi/ocean/pkg/tools/machine"
	"path/filepath"

	"mobingi/ocean/pkg/config"
	"mobingi/ocean/pkg/constants"
	cmdutil "mobingi/ocean/pkg/util/cmd"
)

const schedulerServiceTemplate = `[Unit]
Description=Kubernetes Scheduler
Documentation=https://github.com/GoogleCloudPlatform/kubernetes
After=network.target
After=kube-apiserver.service

[Service]
ExecStart=/usr/local/bin/kube-scheduler \
  --bind-address=127.0.0.1 \
  --leader-elect=true \
  --kubeconfig=/etc/kubernetes/scheduler.conf \
  --authentication-kubeconfig=/etc/kubernetes/scheduler.conf \
  --authorization-kubeconfig=/etc/kubernetes/scheduler.conf
Restart=on-failure
RestartSec=5

[Install]
WantedBy=multi-user.target`

func NewRunSchedulerJob(cfg *config.Config) (*machine.Job, error) {
	job := machine.NewJob("scheduler-service")

	job.AddCmd(cmdutil.NewWriteCmd(filepath.Join(constants.ServiceDir, constants.KubeSchedulerService), schedulerServiceTemplate))
	job.AddCmd(cmdutil.NewSystemStartCmd(constants.KubeSchedulerService))

	return job, nil
}