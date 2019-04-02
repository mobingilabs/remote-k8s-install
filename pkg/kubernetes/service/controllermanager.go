package service

import (
	"path/filepath"

	"mobingi/ocean/pkg/config"
	"mobingi/ocean/pkg/constants"
	"mobingi/ocean/pkg/tools/machine"
	cmdutil "mobingi/ocean/pkg/util/cmd"
)

const controllerManagerServiceTemplate = `[Unit]
Description=Kubernetes Controller Manager
Documentation=https://github.com/GoogleCloudPlatform/kubernetes
After=network.target
After=kube-apiserver.service

[Service]
ExecStart=/usr/local/bin/kube-controller-manager \
  --bind-address=127.0.0.1 \
  --leader-elect=true \
  --kubeconfig=/etc/kubernetes/controller-manager.conf \
  --authentication-kubeconfig=/etc/kubernetes/controller-manager.conf \
  --authorization-kubeconfig=/etc/kubernetes/controller-manager.conf \
  --client-ca-file=/etc/kubernetes/pki/ca.crt \
  --requestheader-client-ca-file=/etc/kubernetes/pki/front-proxy-ca.crt \
  --root-ca-file=/etc/kubernetes/pki/ca.crt \
  --service-account-private-key-file=/etc/kubernetes/pki/sa.key \
  --cluster-signing-cert-file=/etc/kubernetes/pki/ca.crt \
  --cluster-signing-key-file=/etc/kubernetes/pki/ca.key \
  --use-service-account-credentials=true \
  --controllers=*,bootstrapsigner,tokencleaner
Restart=on-failure
RestartSec=5

[Install]
WantedBy=multi-user.target`

func NewRunControllerManagerJob(cfg *config.Config) (*machine.Job, error) {
	job := machine.NewJob("controller-manager-service")

	job.AddCmd(cmdutil.NewWriteCmd(filepath.Join(constants.ServiceDir, constants.KubeControllerManagerService), controllerManagerServiceTemplate))
	job.AddCmd(cmdutil.NewSystemStartCmd(constants.KubeControllerManagerService))

	return job, nil
}
