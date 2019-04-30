package app

import (
	"context"
	pb "mobingi/ocean/app/proto"
	"mobingi/ocean/pkg/config"
	"mobingi/ocean/pkg/log"
	"mobingi/ocean/pkg/phases"
	"mobingi/ocean/pkg/storage"
	"time"
)

type cluster struct{}

/**
TODO:
	cluster name write to database, stop if exist
	remove master and node wait for optimize
	kubelet wait for optimize
	how to download kubelet
	yum install config use aliconfig
*/
func (c *cluster) Init(ctx context.Context, clusterCfg *pb.ClusterConfig) (*pb.Response, error) {
	// Get cluster config
	cfg, err := config.LoadConfigFromGrpc(clusterCfg)
	if err != nil {
		return nil, err
	}
	// TODO remove all documents
	drop := storage.ClusterMongo{}
	drop.RemoveAllDocuments()
	// Create certs and kubeconfs
	storage, err := phases.Init(cfg)
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
	// Store bootstrap conf in the database
	// TODO waiting for apiserver to start
	time.Sleep(10 * time.Second)
	err = storage.SetBootstrapConf(cfg.ClusterName)
	if err != nil {
		return nil, err
	}

	return &pb.Response{Message: ""}, nil
}
