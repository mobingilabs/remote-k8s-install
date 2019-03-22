package kubescheduler

const serviceTemplate = `[Unit]
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
