package app

import (
	"context"
	pb "mobingi/ocean/app/proto"
	"mobingi/ocean/pkg/kubernetes/client"
	"mobingi/ocean/pkg/storage"
)

type node struct{}

func (n *node) Join(ctx context.Context, cfg *pb.InstanceNode) (*pb.NodeConfs, error) {
	// TODO wait for auto choose cluster
	var clusterName = "kubernetes"

	storage := storage.NewStorage()
	bootstrapconf, err := storage.GetKubeconf(clusterName, "bootstrap.conf")
	if err != nil {
		return nil, err
	}

	var certs []*pb.Cert
	certsMap, err := storage.AllCerts(clusterName)
	if err != nil {
		return nil, err
	}
	for name, cert := range certsMap {
		certs = append(certs, &pb.Cert{
			Name: name,
			Cert: cert,
		})
	}

	return &pb.NodeConfs{
		BootstrapConf: bootstrapconf,
		Certs:         certs,
	}, nil
}

func (n *node) Delete(ctx context.Context, cfg *pb.InstanceNode) (*pb.Response, error) {

	storage := storage.NewStorage()
	kubeconfig, err := storage.GetKubeconf("kubernetes", "admin.conf")
	if err != nil {
		return nil, err
	}
	err = client.Init(kubeconfig)
	if err != nil {
		return nil, err
	}

	err = client.DeleteNode(cfg.InstanceID)
	if err != nil {
		return nil, err
	}

	return &pb.Response{Message: ""}, nil
}
