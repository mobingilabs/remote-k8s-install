package staticpod

const controllerManagerYaml = `
apiVersion: v1
kind: Pod
metadata:
  annotations:
    scheduler.alpha.kubernetes.io/critical-pod: ""
  creationTimestamp: null
  labels:
    component: kube-controller-manager
    tier: control-plane
  name: kube-controller-manager
  namespace: kube-system
spec:
  containers:
  - command:
    - kube-controller-manager
    - --bind-address=127.0.0.1
    - --leader-elect=true
    - --allocate-node-cidrs=true
    - --cluster-cidr=10.10.0.1/24
    - --kubeconfig=/etc/kubernetes/controller-manager.conf
    - --authentication-kubeconfig=/etc/kubernetes/controller-manager.conf
    - --authorization-kubeconfig=/etc/kubernetes/controller-manager.conf
    - --client-ca-file=/etc/kubernetes/pki/ca.crt
    - --requestheader-client-ca-file=/etc/kubernetes/pki/front-proxy-ca.crt
    - --root-ca-file=/etc/kubernetes/pki/ca.crt
    - --service-account-private-key-file=/etc/kubernetes/pki/sa.key
    - --cluster-signing-cert-file=/etc/kubernetes/pki/ca.crt
    - --cluster-signing-key-file=/etc/kubernetes/pki/ca.key
    - --use-service-account-credentials=true
    - --controllers=*,bootstrapsigner,tokencleaner
    image: cnbailian/kube-controller-manager:v1.13.3
    name: kube-controller-manager
    volumeMounts:
    - mountPath: /etc/kubernetes/
      name: k8s-confs
      readOnly: true
  hostNetwork: true
  volumes:
  - hostPath:
      path: /etc/kubernetes/
      type: DirectoryOrCreate
    name: k8s-confs
status: {}`
