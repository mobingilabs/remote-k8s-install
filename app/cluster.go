package app

import (
	"context"
	"fmt"
	pb "mobingi/ocean/app/proto"
	"mobingi/ocean/pkg/config"
	"mobingi/ocean/pkg/kubernetes/client"
	"mobingi/ocean/pkg/kubernetes/client/nodes"
	"mobingi/ocean/pkg/kubernetes/service"
	"mobingi/ocean/pkg/log"
	"mobingi/ocean/pkg/phases"
	"mobingi/ocean/pkg/services/tencent"
	configstorage "mobingi/ocean/pkg/storage"
)

type cluster struct{}

func (c *cluster) Init(ctx context.Context, clusterCfg *pb.ClusterConfig) (*pb.Response, error) {
	// Get cluster config
	cfg, err := config.LoadConfigFromGrpc(clusterCfg)
	if err != nil {
		return nil, err
	}
	// Cluster exist validate
	storage := configstorage.NewStorage()
	exist, err := storage.Exist(cfg.ClusterName)
	if err != nil {
		return nil, err
	}
	if exist {
		return nil, fmt.Errorf("Cluster: %s already exists", cfg.ClusterName)
	}
	// Create certs and kubeconfs
	_, err = phases.Init(cfg)
	if err != nil {
		return nil, err
	}
	log.Info("Init configs")
	// Init master server
	master := master{}
	for _, v := range cfg.Masters {
		masterCfg := &pb.ServerConfig{
			ClusterName: cfg.ClusterName,
			PublicIP:    v.PublicIP,
			PrivateIP:   v.PrivateIP,
			User:        v.User,
			Password:    v.Password,
		}
		_, err := master.Join(context.Background(), masterCfg)
		if err != nil {
			return nil, err
		}
	}
	log.Info("Cluster initialized")

	adminconf, err := storage.GetKubeconf(cfg.ClusterName, "admin.conf")
	if err != nil {
		return nil, err
	}
	err = service.WaitingForApiserverStart(adminconf)
	if err != nil {
		return nil, err
	}
	log.Info("Apiserver started")

	log.Info("Create node instance")
	insClient := &tencent.InstanceTencent{}
	res, err := insClient.CreateInstance(clusterCfg.NodeNumber)
	if err != nil {
		return nil, err
	}
	nodes.AddNodeFromInstanceIdSet(res, cfg.ClusterName)

	err = client.CreateMonitor(cfg.ClusterName)
	if err != nil {
		return nil, err
	}

	return &pb.Response{Message: ""}, nil
}

func (c *cluster) Delete(ctx context.Context, clusterCfg *pb.ClusterConfig) (*pb.Response, error) {
	// Get cluster config
	cfg, err := config.LoadConfigFromGrpc(clusterCfg)
	if err != nil {
		return nil, err
	}
	// Stop cluster monitors
	err = client.StopMonitor(cfg.ClusterName)
	if err != nil {
		log.Infof("%s 未创建集群监视器", cfg.ClusterName)
	}
	// Delete cluster storage
	storage := configstorage.NewStorage()
	storage.Drop(cfg)
	// Delete masters
	master := master{}
	for _, v := range clusterCfg.Masters {
		_, err = master.Delete(context.Background(), v)
		if err != nil {
			return nil, err
		}
	}
	return &pb.Response{Message: ""}, nil
}
