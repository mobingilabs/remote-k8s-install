package mainfests

import (
	"k8s.io/api/core/v1"
)

func getVolumesController() []v1.Volume {
	return []v1.Volume{newVolume("etc", "/etc/kubernetes")}
}

func getVolumesEtcd() []v1.Volume {
	return []v1.Volume{newVolume("data", "/var/lib/etcd"), newVolume("certs", "/etc/kubernetes/pki/etcd")}
}

func getVolumeMountsController() []v1.VolumeMount {
	return []v1.VolumeMount{newVolumeMounts("etc", "/etc/kubernetes")}
}

func getVolumeMountsEtcd() []v1.VolumeMount {
	return []v1.VolumeMount{
		newVolumeMounts("data", "/var/lib/etcd"),
		newVolumeMounts("certs", "/etc/kubernetes/pki/etcd"),
	}
}

func newVolume(name, path string) v1.Volume {
	pathType := v1.HostPathDirectory
	return v1.Volume{
		Name: name,
		VolumeSource: v1.VolumeSource{
			HostPath: &v1.HostPathVolumeSource{
				Path: path,
				Type: &pathType,
			},
		},
	}
}

func newVolumeMounts(name, mountPath string) v1.VolumeMount {
	return v1.VolumeMount{
		Name:      name,
		MountPath: mountPath,
	}
}
