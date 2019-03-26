systemctl stop kubelet.service
rm /var/lib/kubelet/config.yaml
rm /etc/systemd/system/kubelet.service
rm -rf /var/lib/kubelet
rm -rf /etc/kubernetes

systemctl stop docker.service
yum remove docker \
  docker-client \
                  docker-client-latest \
                  docker-common \
                  docker-latest \
                  docker-latest-logrotate \
                  docker-logrotate \
                  docker-engine
rm -rf /etc/docker