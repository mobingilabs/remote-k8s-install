package app

import (
	"context"
	pb "mobingi/ocean/app/proto"
	"mobingi/ocean/pkg/config"
	"mobingi/ocean/pkg/log"
	"mobingi/ocean/pkg/phases"
)

type cluster struct{}

// TODO cluster name write to database, stop if exist
func (c *cluster) Init(ctx context.Context, clusterCfg *pb.ClusterConfig) (*pb.Response, error) {
	// Get cluster config
	cfg, err := config.LoadConfigFromGrpc(clusterCfg)
	if err != nil {
		return nil, err
	}
	// Create certs and kubeconfs
	err = phases.Init(cfg)
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

	return &pb.Response{Message: ""}, nil
}
