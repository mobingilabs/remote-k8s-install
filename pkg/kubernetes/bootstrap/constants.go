package bootstrap

const (
	// NodeBootstrapperClusterRoleName defines the name of the auto-bootstrapped ClusterRole for letting someone post a CSR
	// TODO: This value should be defined in an other, generic authz package instead of here
	NodeBootstrapperClusterRoleName = "system:node-bootstrapper"
	// NodeKubeletBootstrap defines the name of the ClusterRoleBinding that lets kubelets post CSRs
	NodeKubeletBootstrap = "ocean:kubelet-bootstrap"

	// CSRAutoApprovalClusterRoleName defines the name of the auto-bootstrapped ClusterRole for making the csrapprover controller auto-approve the CSR
	// TODO: This value should be defined in an other, generic authz package instead of here
	// Starting from v1.8, CSRAutoApprovalClusterRoleName is automatically created by the API server on startup
	CSRAutoApprovalClusterRoleName = "system:certificates.k8s.io:certificatesigningrequests:nodeclient"
	// NodeSelfCSRAutoApprovalClusterRoleName is a role defined in default 1.8 RBAC policies for automatic CSR approvals for automatically rotated node certificates
	NodeSelfCSRAutoApprovalClusterRoleName = "system:certificates.k8s.io:certificatesigningrequests:selfnodeclient"
	// NodeAutoApproveBootstrapClusterRoleBinding defines the name of the ClusterRoleBinding that makes the csrapprover approve node CSRs
	NodeAutoApproveBootstrapClusterRoleBinding = "ocean:node-autoapprove-bootstrap"
	// NodeAutoApproveCertificateRotationClusterRoleBinding defines name of the ClusterRoleBinding that makes the csrapprover approve node auto rotated CSRs
	NodeAutoApproveCertificateRotationClusterRoleBinding = "ocean:node-autoapprove-certificate-rotation"

	// BootstrapSignerClusterRoleName sets the name for the ClusterRole that allows access to ConfigMaps in the kube-public ns
	BootstrapSignerClusterRoleName = "Ocean:bootstrap-signer-clusterinfo"
)
