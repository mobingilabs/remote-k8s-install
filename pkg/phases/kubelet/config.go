package kubelet

import (
	kubeletconfigv1beta1 "k8s.io/kubelet/config/v1beta1"
)

// TODO
func newKubeletConfig() *kubeletconfigv1beta1.KubeletConfiguration {
	return &kubeletconfigv1beta1.KubeletConfiguration{
		Address: "0.0.0.0",
	}
}
