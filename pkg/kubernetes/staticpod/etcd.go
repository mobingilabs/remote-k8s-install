package staticpod

import templateutil "mobingi/ocean/pkg/util/template"

const etcdYaml = `
apiVersion: v1
kind: Pod
metadata:
  annotations:
    scheduler.alpha.kubernetes.io/critical-pod: ""
  creationTimestamp: null
  labels:
    component: etcd
    tier: control-plane
  name: etcd
  namespace: kube-system
spec:
  containers:
  - command:
    - etcd
    - --name=etcd1
    - --initial-advertise-peer-urls=https://{{.IP}}:2380
    - --listen-peer-urls=https://{{.IP}}:2380
    - --listen-client-urls=https://{{.IP}}:2379
    - --advertise-client-urls=https://{{.IP}}:2379
    - --initial-cluster-token=etcd-cluster
    - --initial-cluster=etcd1=https://{{.IP}}:2380
    - --initial-cluster-state=new
    - --data-dir=/var/lib/etcd
    - --cert-file=/etc/kubernetes/pki/etcd/server.crt
    - --key-file=/etc/kubernetes/pki/etcd/server.key
    - --client-cert-auth=true
    - --trusted-ca-file=/etc/kubernetes/pki/etcd/ca.crt
    - --peer-cert-file=/etc/kubernetes/pki/etcd/peer.crt
    - --peer-key-file=/etc/kubernetes/pki/etcd/peer.key
    - --peer-client-cert-auth=true
    - --peer-trusted-ca-file=/etc/kubernetes/pki/etcd/ca.crt
    - --peer-key-file=/etc/kubernetes/pki/etcd/peer.key
    image: cnbailian/etcd:3.3.10
    name: etcd
    volumeMounts:
    - mountPath: /etc/kubernetes/pki
      name: k8s-certs
      readOnly: true
  hostNetwork: true
  volumes:
  - hostPath:
      path: /etc/kubernetes/pki
      type: DirectoryOrCreate
    name: k8s-certs
status: {}`

type etcdTemplateData struct {
	IP string
}

func getEtcdStaticPodFile(ip string) ([]byte, error) {
	templateData := etcdTemplateData{
		IP: ip,
	}
	data, err := templateutil.Parse(etcdYaml, templateData)
	if err != nil {
		return nil, err
	}

	return data, nil
}
