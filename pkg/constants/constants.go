package constants

const (
	// env
	// ServiceDir is a dir for systemd service file
	ServiceDir  = "/etc/systemd/system"
	WorkDir     = "/etc/kubernetes"
	PKIDir      = "/etc/kubernetes/pki"
	ETCDDataDir = "/var/lib/etcd"
	BinDir      = "/usr/local/bin"

	// certs
	CACertAndKeyBaseName = "ca"
	CACertCommonName     = "kubernetes"

	APIServerCertAndKeyBaseName = "apiserver"
	APIServerCertCommonName     = "kube-apiserver"

	APIServerKubeletClientCertAndKeyBaseName = "apiserver-kubelet-client"
	APIServerKubeletClientCertCommonName     = "kube-apiserver-kubelet-client"

	FrontProxyCACertAndKeyBaseName = "front-proxy-ca"
	FrontProxyCACertCommonName     = "front-proxy-ca"

	FrontProxyClientCertAndKeyBaseName = "front-proxy-client"
	FrontProxyClientCertCommonName     = "front-proxy-client"

	// key
	ServiceAccountKeyBaseName = "sa"

	// service
	EtcdService                  = "etcd.service"
	KubeApiserverService         = "kube-apiserver.service"
	KubeControllerManagerService = "kube-controller-manager.service"
	KubeSchedulerService         = "kube-scheduler.service"
	KubeletService               = "kubelet.service"

	// k8s default group,policy rule has been created
	MastersGroup = "system:masters"
	NodesGroup   = "system:nodes"

	NodeBootstrapTokenAuthGroup = "system:bootstrappers:ocean:default-node-token"

	// bootstrap
	NodeKubeletBootstrap                                 = "ocean:kubelet-bootstrap"
	NodeBootstrapperClusterRoleName                      = "system:node-bootstrapper"
	NodeSelfCSRAutoApprovalClusterRoleName               = "system:certificates.k8s.io:certificatesigningrequests:selfnodeclient"
	NodeAutoApproveBootstrapClusterRoleBinding           = "ocean:node-autoapprove-bootstrap"
	CSRAutoApprovalClusterRoleName                       = "system:certificates.k8s.io:certificatesigningrequests:nodeclient"
	NodeAutoApproveCertificateRotationClusterRoleBinding = "ocean:node-autoapprove-certificate-rotation"
)
