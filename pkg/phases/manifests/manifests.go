package manifests

import (
	"encoding/json"
	"fmt"
	"strings"

	"gopkg.in/yaml.v2"
	"k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/intstr"
)

const (
	Etcd                  = "etcd"
	KubeAPIServer         = "kube-apiserver"
	KubeControllerManager = "kube-controller-manager"
	KubeScheduler         = "kube-scheduler"

	etcdClusterToken = "token"
)

type Options struct {
	// now it's machine's private ip
	IPs []string

	EtcdImage              string
	APIServerImage         string
	ControllerManagerImage string
	SchedulerImage         string

	ServiceIPRange string
}

// NewStaticPodManifests will return etcds, apiservers, controllerManager, scheduler
// these map key is private ip
func NewStaticPodManifests(o Options) (map[string][]byte, map[string][]byte, []byte, []byte) {
	etcdPods := make(map[string]*v1.Pod, len(o.IPs))
	etcdArguments := o.createEtcdArguments()
	for k, v := range etcdArguments {
		etcdPods[k] = getEtcdPod(o.EtcdImage, v)
	}
	fillPods(etcdPods, Etcd)
	etcdManifests := make(map[string][]byte, len(o.IPs))
	for k, v := range etcdPods {
		data := marshalToYAML(v)
		etcdManifests[k] = data
	}

	apiServerPods := make(map[string]*v1.Pod, len(o.IPs))
	apiServerArguments := o.createAPIServerArguments()
	for k, v := range apiServerArguments {
		apiServerPods[k] = getAPIServerPod(o.APIServerImage, v)
	}
	// TODO Etcd from external pkg constants
	fillPods(apiServerPods, KubeAPIServer)
	apiServerManifests := make(map[string][]byte, len(o.IPs))
	for k, v := range etcdPods {
		data := marshalToYAML(v)
		apiServerManifests[k] = data
	}

	controllerManagerPod := map[string]*v1.Pod{
		KubeControllerManager: getControllerManagerPod(o.ControllerManagerImage),
	}
	fillPods(controllerManagerPod, KubeControllerManager)

	schedulerPod := map[string]*v1.Pod{
		KubeScheduler: getSchedulerPod(o.SchedulerImage),
	}
	fillPods(schedulerPod, KubeScheduler)

	return etcdManifests, apiServerManifests, marshalToYAML(controllerManagerPod[KubeControllerManager]), marshalToYAML(schedulerPod[KubeScheduler])
}

func (o Options) createEtcdArguments() map[string]map[string]string {
	data := make(map[string]map[string]string)
	nameAndAddress := make([]string, 0, len(o.IPs))
	for i, v := range o.IPs {
		nameAndAddress = append(nameAndAddress, fmt.Sprintf("%s=%s", getNodeName(i), getEtcdPeerURL(v)))
	}
	for i, v := range o.IPs {
		data[v] = map[string]string{
			"name":                        getNodeName(i),
			"initial-advertise-peer-urls": getEtcdPeerURL(v),
			"listen-peer-urls":            getEtcdPeerURL(v),
			"listen-client-urls":          strings.Join([]string{getEtcdClientURL(v), getEtcdClientURL("127.0.0.1")}, ","),
			"advertise-client-urls":       getEtcdClientURL(v),
			"initial-cluster-token":       etcdClusterToken,
			"initial-cluster":             strings.Join(nameAndAddress, ","),
		}
	}

	return data
}

func (o Options) createAPIServerArguments() map[string]map[string]string {
	etcdServers := make([]string, 0, len(o.IPs))
	for _, v := range o.IPs {
		etcdServers = append(etcdServers, getEtcdClientURL(v))
	}

	data := make(map[string]map[string]string, len(o.IPs))
	for _, v := range o.IPs {
		data[v] = map[string]string{
			"advertise-address":        v,
			"etcd-servers":             strings.Join(etcdServers, ","),
			"service-cluster-ip-range": o.ServiceIPRange,
		}
	}

	return data
}

func getEtcdPod(image string, arguments map[string]string) *v1.Pod {
	return &v1.Pod{
		Spec: v1.PodSpec{
			Containers: []v1.Container{
				{
					Image:        image,
					Command:      getEtcdCommand(arguments),
					VolumeMounts: getVolumeMountsEtcd(),
				},
			},
			Volumes: getVolumesEtcd(),
		},
	}
}

// TODO use constans replace some arguments
func getEtcdCommand(extraArguments map[string]string) []string {
	defaultArguments := map[string]string{
		"data-dir":              "/var/lib/etcd",
		"initial-cluster-state": "new",
		"client-cert-auth":      "true",
		"cert-file":             "/etc/kubernetes/pki/etcd-server.crt",
		"key-file":              "/etc/kubernetes/pki/etcd-server.key",
		"trusted-ca-file":       "/etc/kubernetes/pki/ca.crt",
		//"peer-client-cert-auth": "true",
		//"peer-trusted-ca-file":  "/etc/kubernetes/pki/etcd/ca.crt",
		//"peer-key-file":         "/etc/kubernetes/pki/etcd/peer.key",
		//"peer-cert-file":        "/etc/kubernetes/pki/etcd/peer.crt",
	}

	return buildCommand(Etcd, defaultArguments, extraArguments)
}

func getAPIServerPod(image string, extraArguments map[string]string) *v1.Pod {
	return &v1.Pod{
		Spec: v1.PodSpec{
			Containers: []v1.Container{
				{
					Image:   image,
					Command: getAPIServerCommand(extraArguments),
					//					LivenessProbe: newLivenessProbe("11", 6443),
					VolumeMounts: getVolumeMountsController(),
				},
			},
			Volumes: getVolumesController(),
		},
	}
}

func getControllerManagerPod(image string) *v1.Pod {
	return &v1.Pod{
		Spec: v1.PodSpec{
			Containers: []v1.Container{
				{
					Image:   image,
					Command: getControllerManagerCommand(nil),
					//				LivenessProbe: newLivenessProbe("11", 10544),
					VolumeMounts: getVolumeMountsController(),
				},
			},
			Volumes: getVolumesController(),
		},
	}
}

func getSchedulerPod(image string) *v1.Pod {
	return &v1.Pod{
		Spec: v1.PodSpec{
			Containers: []v1.Container{
				{
					Image:   image,
					Command: getSchedulerCommand(nil),
					//			LivenessProbe: newLivenessProbe("11", 10533),
					VolumeMounts: getVolumeMountsController(),
				},
			},
			Volumes: getVolumesController(),
		},
	}
}

func fillPods(pods map[string]*v1.Pod, name string) {
	for _, v := range pods {
		v.APIVersion = v1.SchemeGroupVersion.String()
		v.Kind = "Pod"
		v.Name = name
		v.Namespace = "kube-system"
		v.Labels = map[string]string{
			"component": name,
			"tier":      "control-plane",
		}

		v.Spec.Containers[0].Name = name
		v.Spec.Containers[0].ImagePullPolicy = v1.PullIfNotPresent
		v.Spec.HostNetwork = true
		v.Spec.PriorityClassName = "system-cluster-critical"
	}
}

func getAPIServerCommand(extraArguments map[string]string) []string {
	defaultArguments := map[string]string{
		"allow-privileged":            "true",
		"authorization-mode":          "Node,RBAC",
		"insecure-port":               "0",
		"secure-port":                 "6443",
		"anoymous-auth":               "false",
		"enable-admission-plugins":    "NodeRestriction",
		"enable-bootstrap-token-auth": "true",
		"service-account-key-file":    "/etc/kubernetes/pki/sa.pub",

		"client-ca-file": "/etc/kubernetes/pki/ca.crt",
		//"etcd-cafile":                "/etc/kubernetes/pki/ca.crt",
		//"etcd-certfile":              "/etc/kubernetes/pki/apiserver-etcd-client.crt",
		//"etcd-keyfile":               "/etc/kubernetes/pki/apiserver-etcd-client.key",
		"kubelet-client-certificate": "/etc/kubernetes/pki/apiserver-kubelet-client.crt",
		"kubelet-client-key":         "/etc/kubernetes/pki/apiserver-kubelet-client.key",
		"tls-cert-file":              "/etc/kubernetes/pki/apiserver.crt",
		"tls-private-key-file":       "/etc/kubernetes/pki/apiserver.key",
	}

	return buildCommand(KubeAPIServer, defaultArguments, extraArguments)
}

func getControllerManagerCommand(extraArguments map[string]string) []string {
	defaultArguments := map[string]string{
		"bind-address":                     "127.0.0.1",
		"kubeconfig":                       "/etc/kubernetes/controller-manager.conf",
		"root-ca-file":                     "/etc/kubernetes/pki/ca.crt",
		"client-ca-file":                   "/etc/kubernetes/pki/ca.crt",
		"cluster-signing-cert-file":        "/etc/kubernetes/pki/ca.crt",
		"cluster-signing-key-file":         "/etc/kubernetes/pki/ca.key",
		"controllers":                      "*,bootstrapsigner,tokencleaner",
		"leader-elect":                     "true",
		"service-account-private-key-file": "/etc/kubernetes/pki/sa.key",
		"use-service-account-credentials":  "true",
	}

	return buildCommand(KubeControllerManager, defaultArguments, extraArguments)
}

func getSchedulerCommand(extraArguments map[string]string) []string {
	defaultArguments := map[string]string{
		"bind-address": "127.0.0.1",
		"kubeconfig":   "/etc/kubernetes/scheduler.conf",
		"leader-elect": "true",
	}

	return buildCommand(KubeScheduler, defaultArguments, extraArguments)
}

func buildCommand(baseCommand string, defaultArguments map[string]string, extraArguments map[string]string) []string {
	command := []string{}
	arguments := map[string]string{}

	for k, v := range defaultArguments {
		arguments[k] = v
	}
	for k, v := range extraArguments {
		arguments[k] = v
	}

	command = append(command, baseCommand)
	for k, v := range arguments {
		command = append(command, fmt.Sprintf("--%s=%s", k, v))
	}

	return command
}

func newLivenessProbe(host string, port int) *v1.Probe {
	return &v1.Probe{
		Handler: v1.Handler{
			HTTPGet: &v1.HTTPGetAction{
				Host: host,
				Path: "/healthz",
				Port: intstr.FromInt(port),
			},
		},
		InitialDelaySeconds: 15,
		TimeoutSeconds:      15,
		FailureThreshold:    8,
	}
}

func marshalToYAML(obj runtime.Object) []byte {
	data, _ := json.Marshal(obj)
	temMap := map[string]interface{}{}
	json.Unmarshal(data, &temMap)
	res, _ := yaml.Marshal(temMap)
	return res
}

// TODO use https
func getEtcdPeerURL(ip string) string {
	return fmt.Sprintf("http://%s:2380", ip)
}

func getEtcdClientURL(ip string) string {
	return fmt.Sprintf("http://%s:2379", ip)
}

func getNodeName(i int) string {
	return fmt.Sprintf("node:%d", i)
}
