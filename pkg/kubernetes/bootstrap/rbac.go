package bootstrap

import (
	rbac "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apiserver/pkg/authentication/user"
	clientset "k8s.io/client-go/kubernetes"
	bootstrapapi "k8s.io/cluster-bootstrap/token/api"
	rbacv1 "k8s.io/kubernetes/pkg/apis/rbac/v1"

	"mobingi/ocean/pkg/constants"
)

// CreateClusterInfoRBACRules creates the RBAC rules for exposing the cluster-info ConfigMap in the kube-public namespace to unauthenticated users
func CreateClusterInfoRBACRules(client clientset.Interface) error {
	role := &rbac.Role{
		ObjectMeta: metav1.ObjectMeta{
			Name:      BootstrapSignerClusterRoleName,
			Namespace: metav1.NamespacePublic,
		},
		Rules: []rbac.PolicyRule{
			rbacv1.NewRule("get").Groups("").Resources("configmaps").Names(bootstrapapi.ConfigMapClusterInfo).RuleOrDie(),
		},
	}
	if _, err := client.RbacV1().Roles(role.ObjectMeta.Namespace).Create(role); err != nil {
		return err
	}

	roleBinding := &rbac.RoleBinding{
		ObjectMeta: metav1.ObjectMeta{
			Name:      BootstrapSignerClusterRoleName,
			Namespace: metav1.NamespacePublic,
		},
		RoleRef: rbac.RoleRef{
			APIGroup: rbac.GroupName,
			Kind:     "Role",
			Name:     BootstrapSignerClusterRoleName,
		},
		Subjects: []rbac.Subject{
			{
				Kind: rbac.UserKind,
				Name: user.Anonymous,
			},
		},
	}
	if _, err := client.RbacV1().RoleBindings(roleBinding.ObjectMeta.Namespace).Create(roleBinding); err != nil {
		return err
	}

	return nil
}

// AllowBootstrapTokensToPostCSRs creates RBAC rules in a way the makes Node Bootstrap Tokens able to post CSRs
func AllowBootstrapTokensToPostCSRs(client clientset.Interface) error {
	roleBinding := &rbac.ClusterRoleBinding{
		ObjectMeta: metav1.ObjectMeta{
			Name: NodeKubeletBootstrap,
		},
		RoleRef: rbac.RoleRef{
			APIGroup: rbac.GroupName,
			Kind:     "ClusterRole",
			Name:     NodeBootstrapperClusterRoleName,
		},
		Subjects: []rbac.Subject{
			{
				Kind: rbac.GroupKind,
				Name: constants.NodeBootstrapTokenAuthGroup,
			},
		},
	}
	if _, err := client.RbacV1().RoleBindings(roleBinding.ObjectMeta.Namespace).Create(roleBinding); err != nil {
		return err
	}

	return nil
}

// AutoApproveNodeBootstrapTokens creates RBAC rules in a way that makes Node Bootstrap Tokens' CSR auto-approved by the csrapprover controller
func AutoApproveNodeBootstrapTokens(client clientset.Interface) error {
	roleBinding := &rbac.ClusterRoleBinding{
		ObjectMeta: metav1.ObjectMeta{
			Name: NodeAutoApproveBootstrapClusterRoleBinding,
		},
		RoleRef: rbac.RoleRef{
			APIGroup: rbac.GroupName,
			Kind:     "ClusterRole",
			Name:     CSRAutoApprovalClusterRoleName,
		},
		Subjects: []rbac.Subject{
			{
				Kind: "Group",
				Name: constants.NodeBootstrapTokenAuthGroup,
			},
		},
	}
	if _, err := client.RbacV1().RoleBindings(roleBinding.ObjectMeta.Namespace).Create(roleBinding); err != nil {
		return err
	}

	return nil
}

// AutoApproveNodeCertificateRotation creates RBAC rules in a way that makes Node certificate rotation CSR auto-approved by the csrapprover controller
func AutoApproveNodeCertificateRotation(client clientset.Interface) error {
	roleBinding := &rbac.ClusterRoleBinding{
		ObjectMeta: metav1.ObjectMeta{
			Name: NodeAutoApproveCertificateRotationClusterRoleBinding,
		},
		RoleRef: rbac.RoleRef{
			APIGroup: rbac.GroupName,
			Kind:     "ClusterRole",
			Name:     NodeSelfCSRAutoApprovalClusterRoleName,
		},
		Subjects: []rbac.Subject{
			{
				Kind: "Group",
				Name: constants.NodesGroup,
			},
		},
	}
	if err := client.RbacV1().RoleBindings(roleBinding.ObjectMeta.Namespace).Create(roleBinding); err != nil {
		return err
	}

	return nil
}
