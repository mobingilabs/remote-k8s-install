systemctl stop kube-controller-manager
systemctl stop kube-scheduler
systemctl stop kube-apiserver
systemctl stop etcd

rm /etc/systemd/system/etcd.service 
rm /etc/systemd/system/kube-apiserver.service 
rm /etc/systemd/system/kube-controller-manager.service 
rm /etc/systemd/system/kube-scheduler.service 

#k8s workdir
rm -rf /etc/kubernetes
#etcd datadir
rm -rf /var/lib/etcd