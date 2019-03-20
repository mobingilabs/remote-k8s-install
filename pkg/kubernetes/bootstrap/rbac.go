package bootstrap

import (
	rbac "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	clientset "k8s.io/client-go/kubernetes"

	"mobingi/ocean/pkg/constants"
)

// AllowBootstrapTokensToPostCSRs creates RBAC rules in a way the makes Node Bootstrap Tokens able to post CSRs
func AllowBootstrapTokensToPostCSRs(client clientset.Interface) error {
	roleBinding := &rbac.ClusterRoleBinding{
		ObjectMeta: metav1.ObjectMeta{
			Name: constants.NodeKubeletBootstrap,
		},
		RoleRef: rbac.RoleRef{
			APIGroup: rbac.GroupName,
			Kind:     "ClusterRole",
			Name:     constants.NodeBootstrapperClusterRoleName,
		},
		Subjects: []rbac.Subject{
			{
				Kind: rbac.GroupKind,
				Name: constants.NodeBootstrapTokenAuthGroup,
			},
		},
	}
	if _, err := client.RbacV1().ClusterRoleBindings().Create(roleBinding); err != nil {
		return err
	}

	return nil
}

// AutoApproveNodeBootstrapTokens creates RBAC rules in a way that makes Node Bootstrap Tokens' CSR auto-approved by the csrapprover controller
func AutoApproveNodeBootstrapTokens(client clientset.Interface) error {
	roleBinding := &rbac.ClusterRoleBinding{
		ObjectMeta: metav1.ObjectMeta{
			Name: constants.NodeAutoApproveBootstrapClusterRoleBinding,
		},
		RoleRef: rbac.RoleRef{
			APIGroup: rbac.GroupName,
			Kind:     "ClusterRole",
			Name:     constants.CSRAutoApprovalClusterRoleName,
		},
		Subjects: []rbac.Subject{
			{
				Kind: "Group",
				Name: constants.NodeBootstrapTokenAuthGroup,
			},
		},
	}
	if _, err := client.RbacV1().ClusterRoleBindings().Create(roleBinding); err != nil {
		return err
	}

	return nil
}

// AutoApproveNodeCertificateRotation creates RBAC rules in a way that makes Node certificate rotation CSR auto-approved by the csrapprover controller
func AutoApproveNodeCertificateRotation(client clientset.Interface) error {
	roleBinding := &rbac.ClusterRoleBinding{
		ObjectMeta: metav1.ObjectMeta{
			Name: constants.NodeAutoApproveCertificateRotationClusterRoleBinding,
		},
		RoleRef: rbac.RoleRef{
			APIGroup: rbac.GroupName,
			Kind:     "ClusterRole",
			Name:     constants.NodeSelfCSRAutoApprovalClusterRoleName,
		},
		Subjects: []rbac.Subject{
			{
				Kind: "Group",
				Name: constants.NodesGroup,
			},
		},
	}
	if _, err := client.RbacV1().ClusterRoleBindings().Create(roleBinding); err != nil {
		return err
	}

	return nil
}
