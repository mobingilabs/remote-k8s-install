package constants

const (
	// ServiceDir is a dir for systemd service file
	ServiceDir = "/etc/systemd/system"

	EtcdService                  = "etcd.service"
	KubeApiserverService         = "kube-apiserver.service"
	KubeControllerManagerService = "kube-controller-manager.service"
	KubeSchedulerService         = "kube-scheduler.service"

	MastersGroup = "system:masters"
	NodesGroup   = "system:nodes"

	NodeBootstrapTokenAuthGroup = "system:bootstrappers:kubeadm:default-node-token"
)
