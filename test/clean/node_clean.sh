systemctl stop kubelet.service
rm /var/lib/kubelet/config.yaml
rm /etc/systemd/system/kubelet.service