package constants

const (
	// ServiceDir is a dir for systemd service file
	ServiceDir = "/etc/systemd/system"

	WorkDir = "/etc/kubernetes"
	PKIDir  = "/etc/kubernetes/pki"

	EtcdService                  = "etcd.service"
	KubeApiserverService         = "kube-apiserver.service"
	KubeControllerManagerService = "kube-controller-manager.service"
	KubeSchedulerService         = "kube-scheduler.service"
	KubeletService               = "kubelet.service"

	// k8s default group,policy rule has been created
	MastersGroup = "system:masters"
	NodesGroup   = "system:nodes"

	NodeBootstrapTokenAuthGroup = "system:bootstrappers:kubeadm:default-node-token"
)
