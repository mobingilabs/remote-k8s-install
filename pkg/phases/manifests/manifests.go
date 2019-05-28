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
)

type Options struct {
	IP        string
	NodeName  string
	EtcdToken string
	// ohter two master's info, k is name,v is IP
	Companions map[string]string

	EtcdImage              string
	APIServerImage         string
	ControllerManagerImage string
	SchedulerImage         string

	ServiceIPRange string
}

// GetStaticPodMainfests create etcd,apiserver,controller-manager,scheduler staticPod files
func NewStaticPodManifests(o Options) map[string][]byte {
	pods := map[string]*v1.Pod{
		Etcd:                  getEtcdPod(o.EtcdImage, o.createEtcdArguments()),
		KubeAPIServer:         getAPIServerPod(o.APIServerImage, o.createAPIServerArguments()),
		KubeControllerManager: getControllerManagerPod(o.ControllerManagerImage),
		KubeScheduler:         getSchedulerPod(o.SchedulerImage),
	}
	fillPods(pods)

	mainfests := make(map[string][]byte, len(pods))
	for k, v := range pods {
		data := marshalToYAML(v)
		mainfests[getFileName(k)] = data
	}

	return mainfests
}

func (o Options) createEtcdArguments() map[string]string {
	initialClusters := make([]string, 0, len(o.Companions)+1)
	for k, v := range o.Companions {
		initialClusters = append(initialClusters, fmt.Sprintf("%s=%s", k, getEtcdPeerURL(v)))
	}
	initialClusters = append(initialClusters, fmt.Sprintf("%s=%s", o.NodeName, getEtcdPeerURL(o.IP)))
	return map[string]string{
		"name":                        o.NodeName,
		"initial-advertise-peer-urls": getEtcdPeerURL(o.IP),
		"listen-peer-urls":            getEtcdPeerURL(o.IP),
		"listen-client-urls":          strings.Join([]string{getEtcdClientURL(o.IP), getEtcdClientURL("127.0.0.1")}, ","),
		"advertise-client-urls":       getEtcdClientURL(o.IP),
		"initial-cluster-token":       o.EtcdToken,
		"initial-cluster":             strings.Join(initialClusters, ","),
	}
}

func (o Options) createAPIServerArguments() map[string]string {
	ips := make([]string, 0, len(o.Companions)+1)
	ips = append(ips, getEtcdClientURL(o.IP))
	for _, v := range o.Companions {
		ips = append(ips, getEtcdClientURL(v))
	}
	return map[string]string{
		"advertise-address":        o.IP,
		"etcd-servers":             strings.Join(ips, ","),
		"service-cluster-ip-range": o.ServiceIPRange,
	}
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

func fillPods(pods map[string]*v1.Pod) {
	for k, v := range pods {
		v.APIVersion = v1.SchemeGroupVersion.String()
		v.Kind = "Pod"
		v.Name = k
		v.Namespace = "kube-system"
		v.Labels = map[string]string{
			"component": k,
			"tier":      "control-plane",
		}

		v.Spec.Containers[0].Name = k
		v.Spec.Containers[0].ImagePullPolicy = v1.PullIfNotPresent
		v.Spec.HostNetwork = true
		v.Spec.PriorityClassName = "system-cluster-critical"
	}
}

func getAPIServerCommand(extraArguments map[string]string) []string {
	defaultArguments := map[string]string{
		"allow-privileged":   "true",
		"authorization-mode": "Node,RBAC",
		"insecure-port":      "0",
		"secure-port":        "6443",
		// 		"anoymous-auth":               "false",
		"enable-admission-plugins":    "NodeRestriction",
		"enable-bootstrap-token-auth": "true",

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
		"bind-address":              "127.0.0.1",
		"kubeconfig":                "/etc/kubernetes/controller-manager.conf",
		"root-ca-file":              "/etc/kubernetes/pki/ca.crt",
		"client-ca-file":            "/etc/kubernetes/pki/ca.crt",
		"cluster-signing-cert-file": "/etc/kubernetes/pki/ca.crt",
		"cluster-signing-key-file":  "/etc/kubernetes/pki/ca.key",
		"controllers":               "*,bootstrapsigner,tokencleaner",
		"leader-elect":              "true",
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

func getFileName(name string) string {
	return fmt.Sprintf("%s.%s", name, "yaml")
}
