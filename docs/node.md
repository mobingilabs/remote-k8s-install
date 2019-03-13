1.systemctl status etcd 显示 active，但是无法使用etcd，建议用etcdctl进行测试
2.最大的坑，systemd的service文件，[Service]下面的Type=notify,勿设置，不然程序超时退出