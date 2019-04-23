package service

import (
	"path/filepath"

	"mobingi/ocean/pkg/constants"
	"mobingi/ocean/pkg/tools/machine"
	cmdutil "mobingi/ocean/pkg/util/cmd"
)

const kubeletServiceTemplate = `[Unit]
Description= The Kubernetes Node Agent
Documentation=https://github.com/kubernetes.io/docs
After=network.target

[Service]
ExecStart=/usr/local/bin/kubelet 
Rstart=on-failure
RestartSec=5

[Install]
WantedBy=multi-user.target`

const kubeletServicedDir = "kubelet.service.d"
const kubeletServicedName = "10-ocean.conf"
const kubeletServicedFileContent = `
[Service]
Environment="KUBELET_KUBECONFIG_ARGS=--bootstrap-kubeconfig=/etc/kubernetes/bootstrap-kubelet.conf --kubeconfig=/etc/kubernetes/kubelet.conf"
Environment="KUBELET_CONFIG_ARGS=--config=/var/lib/kubelet/config.yaml"
EnvironmentFile=-/var/lib/kubelet/ocean-flags.env
EnvironmentFile=-/etc/sysconfig/kubelet
ExecStart=
ExecStart=/usr/local/bin/kubelet \$KUBELET_KUBECONFIG_ARGS \$KUBELET_CONFIG_ARGS --allow-privileged=true
`

// var/lib/kubelet/config.yaml
// TOOD read from config
const kubeletConfigDir = "/var/lib/kubelet"
const kubeletConfigName = "config.yaml"
const kubeletConfigYAML = `
address: 0.0.0.0
apiVersion: kubelet.config.k8s.io/v1beta1
authentication:
  anonymous:
    enabled: false
  webhook:
    cacheTTL: 2m0s
    enabled: true
  x509:
    clientCAFile: /etc/kubernetes/pki/ca.crt
authorization:
  mode: Webhook
  webhook:
    cacheAuthorizedTTL: 5m0s
    cacheUnauthorizedTTL: 30s
cgroupDriver: systemd
cgroupsPerQOS: true
clusterDNS:
- 10.96.0.10
clusterDomain: cluster.local
configMapAndSecretChangeDetectionStrategy: Watch
containerLogMaxFiles: 5
containerLogMaxSize: 10Mi
contentType: application/vnd.kubernetes.protobuf
cpuCFSQuota: true
cpuCFSQuotaPeriod: 100ms
cpuManagerPolicy: none
cpuManagerReconcilePeriod: 10s
enableControllerAttachDetach: true
enableDebuggingHandlers: true
enforceNodeAllocatable:
- pods
eventBurst: 10
eventRecordQPS: 5
evictionHard:
  imagefs.available: 15%
  memory.available: 100Mi
  nodefs.available: 10%
  nodefs.inodesFree: 5%
evictionPressureTransitionPeriod: 5m0s
failSwapOn: true
fileCheckFrequency: 20s
hairpinMode: promiscuous-bridge
healthzBindAddress: 127.0.0.1
healthzPort: 10248
httpCheckFrequency: 20s
imageGCHighThresholdPercent: 85
imageGCLowThresholdPercent: 80
imageMinimumGCAge: 2m0s
iptablesDropBit: 15
iptablesMasqueradeBit: 14
kind: KubeletConfiguration
kubeAPIBurst: 10
kubeAPIQPS: 5
makeIPTablesUtilChains: true
maxOpenFiles: 1000000
maxPods: 110
nodeLeaseDurationSeconds: 40
nodeStatusReportFrequency: 1m0s
nodeStatusUpdateFrequency: 10s
oomScoreAdj: -999
podPidsLimit: -1
port: 10250
registryBurst: 10
registryPullQPS: 5
resolvConf: /etc/resolv.conf
rotateCertificates: true
runtimeRequestTimeout: 2m0s
serializeImagePulls: true
streamingConnectionIdleTimeout: 4h0m0s
syncFrequency: 1m0s
volumeStatsAggPeriod: 1m0s`

const kubeletFlagsFileName = "ocean-flags.env"
const kubeletFlagsContent = `KUBELET_KUBEADM_ARGS=--cgroup-driver=systemd --network-plugin=cni --pod-infra-container-image=cnbailian/pause:3.1`

func NewRunKubeletJob() *machine.Job {
	job := machine.NewJob("kubelet-service")

	job.AddCmd(cmdutil.NewMkdirAllCmd(filepath.Join(constants.ServiceDir, kubeletServicedDir)))
	job.AddCmd(cmdutil.NewMkdirAllCmd(kubeletConfigDir))
	job.AddCmd(cmdutil.NewWriteCmd(filepath.Join(constants.ServiceDir, constants.KubeletService), kubeletServiceTemplate))
	job.AddCmd(cmdutil.NewWriteCmd(filepath.Join(constants.ServiceDir, kubeletServicedDir, kubeletServicedName), kubeletServicedFileContent))
	job.AddCmd(cmdutil.NewWriteCmd(filepath.Join(kubeletConfigDir, kubeletConfigName), kubeletConfigYAML))
	job.AddCmd(cmdutil.NewWriteCmd(filepath.Join(kubeletConfigDir, kubeletFlagsFileName), kubeletFlagsContent))
	job.AddCmd(cmdutil.NewSystemStartCmd(constants.KubeletService))

	return job
}
