package master

import (
	"mobingi/ocean/pkg/constants"
	"mobingi/ocean/pkg/tools/machine"
	cmdutil "mobingi/ocean/pkg/util/cmd"
	"path/filepath"

	_ "github.com/go-sql-driver/mysql"
)

func WriteBootstrapconf(j *machine.Job, bootstrapconf []byte) {
	j.AddCmd(cmdutil.NewWriteCmd(filepath.Join(constants.WorkDir, constants.BootstrapKubeletConfName), string(bootstrapconf)))
}

func InstallDocker(j *machine.Job) {
	j.AddCmd("yum install -y yum-utils device-mapper-perisitent-data lvm2")
	j.AddCmd("yum-config-manager --add-repo https://download.docker.com/linux/centos/docker-ce.repo")
	j.AddCmd("yum install -y docker-ce-18.06.0.ce-3.el7")
	j.AddCmd(cmdutil.NewMkdirAllCmd("/etc/docker"))

	dockerDaemonJson := `{
	"exec-opts": ["native.cgroupdriver=systemd"],
	"log-driver": "json-file",
	"log-opts": {
		"max-size": "100m"
	},
	"storage-driver": "overlay2"
}`
	dockerDaemonJsonName := "daemon.json"
	j.AddCmd(cmdutil.NewWriteCmd(filepath.Join("/etc/docker", dockerDaemonJsonName), dockerDaemonJson))
	// TODO add check
	j.AddCmd(cmdutil.NewSystemStartCmd("docker"))
}

// This will be a http handler
// func InstallMasters(cfg *config.Config) error {
// 	sans := cfg.GetSANs()
// 	certList, err := certs.CreatePKIAssets(cfg.AdvertiseAddress, cfg.PublicIP, sans)
// 	if err != nil {
// 		log.Panicf("cert create:%s", err.Error())
// 	}
// 	log.Info("cert create")

// 	caCert, caKey, err := getCaCertAndKey(certList)
// 	if err != nil {
// 		log.Panicf("get ca cert and key :%s", err.Error())
// 	}
// 	kubeconfs, err := kubeconf.CreateKubeconf(cfg, caCert, caKey)
// 	if err != nil {
// 		log.Panicf("create kube conf :%s", err.Error())
// 	}
// 	log.Info("kubeconf create")
// 	// TODO we will put confs to store, not cache
// 	cache.Put(constants.KubeconfPrefix, "admin.conf", kubeconfs["admin.conf"])

// 	machines := newMachines(cfg)
// 	log.Info("machine init")

// 	g := group.NewGroup(len(cfg.Masters))
// 	job := preparemaster.NewJob(cfg.DownloadBinSite, certList, kubeconfs)
// 	for _, v := range machines {
// 		m := v
// 		g.Add(func() error {
// 			return m.Run(job)
// 		})
// 	}
// 	errs := g.Run()
// 	for _, v := range errs {
// 		if v != nil {
// 			log.Panicf("master prepare:%s", v.Error())
// 		}
// 	}
// 	log.Info("master prepare")

// 	privateIPs := cfg.GetMasterPrivateIPs()
// 	runEtcdCluster(machines, privateIPs)

// 	etcdServers := service.GetEtcdServers(privateIPs)
// 	runControlPlane(machines, privateIPs, etcdServers, cfg.AdvertiseAddress)

// 	// TODO wait for services up
// 	time.Sleep(time.Second)

// 	return nil
// }

// func runEtcdCluster(machines []machine.Machine, privateIPs []string) {
// 	etcdRunJobs, err := service.NewRunEtcdJobs(privateIPs, CertList)
// 	if err != nil {
// 		panic(err)
// 	}

// 	g := group.NewGroup(len(machines))
// 	for i, v := range machines {
// 		m := v
// 		j := i
// 		g.Add(func() error {
// 			return m.Run(etcdRunJobs[j])
// 		})
// 	}
// 	// TODO we will design a error list type for check easily
// 	errs := g.Run()
// 	for _, v := range errs {
// 		if v != nil {
// 			log.Panicf("etcd run:%s", v.Error())
// 		}
// 	}
// 	log.Info("etcd run")
// }

// func runControlPlane(machines []machine.Machine, privateIPs []string, etcdServers, advertiseAddress string) {
// 	controlPlaneJobs, err := service.NewRunControlPlaneJobs(privateIPs, etcdServers, advertiseAddress)
// 	if err != nil {
// 		log.Panic(err)
// 	}

// 	g := group.NewGroup(len(machines))
// 	for i, v := range machines {
// 		m := v
// 		j := i
// 		g.Add(func() error {
// 			return m.Run(controlPlaneJobs[j])
// 		})
// 	}
// 	errs := g.Run()
// 	for _, v := range errs {
// 		if v != nil {
// 			log.Panicf("control plane:%s", v.Error())
// 		}
// 	}
// 	log.Info("control plane")
// }
