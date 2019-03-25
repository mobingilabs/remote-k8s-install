systemctl stop kube-controller-manager
systemctl stop kube-scheduler
systemctl stop kube-apiserver
systemctl stop etcd

rm /etc/systemd/system/etcd.service 
rm /etc/systemd/system/kube-apiserver.service 
rm /etc/systemd/system/kube-controller-manager.service 
rm /etc/systemd/system/kube-scheduler.service 

rm -rf /etc/kubernetes
rm -rf /var/lib/etcd

rm /etc/sysctl.d/k8s.conf
rm /tmp/master.tgz
rm /usr/local/bin/kube-apiserver
rm /usr/local/bin/kube-controller-manager
rm /usr/local/bin/kube-scheduler