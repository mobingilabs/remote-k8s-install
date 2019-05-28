package kubelet

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	kubeletconfigv1beta1 "k8s.io/kubelet/config/v1beta1"
)

// TODO
func newKubeletConfig() *kubeletconfigv1beta1.KubeletConfiguration {
	return &kubeletconfigv1beta1.KubeletConfiguration{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "kubelet.config.k8s.io/v1beta1",
			Kind:       "KubeletConfiguration",
		},
		Address:       "0.0.0.0",
		StaticPodPath: "/etc/kubernetes/manifestes",
	}
}
