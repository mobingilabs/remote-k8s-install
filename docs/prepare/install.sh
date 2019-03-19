#!/usr/bin/env bash
BIN_DIR=/usr/local/bin
#env set
swapoff -a
cat <<EOF | tee /etc/sysctl.d/k8s.conf
net.ipv4.ip_forward = 1
EOF
sysctl -p /etc/sysctl.d/k8s.conf
############download##############################
#k8s binary
curl -L https://dl.k8s.io/v1.13.3/kubernetes-server-linux-amd64.tar.gz -o /tmp/k8s.tar.gz
tar xzvf /tmp/k8s.tar.gz -C /tmp
mv /tmp/kubernetes/server/bin/kubectl ${BIN_DIR}
mv /tmp/kubernetes/server/bin/kube-apiserver ${BIN_DIR}
mv /tmp/kubernetes/server/bin/kube-controller-manager ${BIN_DIR}
mv /tmp/kubernetes/server/bin/kube-scheduler ${BIN_DIR}
rm -rf /tmp/kubernetes
#CNI
mkdir -p /opt/cni/bin
curl -L https://github.com/containernetworking/plugins/releases/download/v0.7.1/cni-plugins-amd64-v0.7.1.tgz -o /tmp/cni.tar.gz
tar zxvf /tmp/cni.tar.gz -C /opt/cni/bin
rm /tmp/cni.tar.gz
#cfssl
curl -L https://pkg.cfssl.org/R1.2/cfssl_linux-amd64 -o ${BIN_DIR}/cfssl
curl -L https://pkg.cfssl.org/R1.2/cfssljson_linux-amd64 -o ${BIN_DIR}/cfssljson
chmod +x ${BIN_DIR}/cfssl ${BIN_DIR}/cfssljson
#etcd
curl -L https://github.com/etcd-io/etcd/releases/download/v3.2.24/etcd-v3.2.24-linux-amd64.tar.gz -o /tmp/etcd.tar.gz
tar xzvf /tmp/etcd.tar.gz -C /tmp
mv /tmp/etcd-v3.2.24-linux-amd64/etcd ${BIN_DIR}
mv /tmp/etcd-v3.2.24-linux-amd64/etcdctl ${BIN_DIR}
rm -f /tmp/etcd.tar.gz
rm -rf /tmp/etcd-v3.2.24-linux-amd64
