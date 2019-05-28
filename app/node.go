package app

import (
	"context"
	pb "mobingi/ocean/app/proto"
	"mobingi/ocean/pkg/kubernetes/bootstrap"
	"mobingi/ocean/pkg/kubernetes/client"
	"mobingi/ocean/pkg/kubernetes/client/nodes"
	"mobingi/ocean/pkg/services/tencent"
	"mobingi/ocean/pkg/storage"
)

type node struct{}

func (n *node) Join(ctx context.Context, cfg *pb.InstanceNode) (*pb.NodeConfs, error) {
	var clusterName = nodes.Nodes[cfg.InstanceID]

	storage := storage.NewStorage()

	adminConf, err := storage.GetKubeconf(clusterName, "admin.conf")
	bootstrapconf, err := bootstrap.Bootstrap(adminConf)
	if err != nil {
		return nil, err
	}
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
	var clusterName = nodes.Nodes[cfg.InstanceID]

	clientset, err := client.NewClient(clusterName)
	if err != nil {
		return nil, err
	}
	var node = &client.Node{
		Client: clientset,
	}
	err = node.DeleteNode(cfg.InstanceID)
	if err != nil {
		return nil, err
	}

	delete(nodes.Nodes, cfg.InstanceID)

	return &pb.Response{Message: ""}, nil
}

func (n *node) SpotInstanceDestroy(ctx context.Context, cfg *pb.InstanceNode) (*pb.Response, error) {
	var clusterName = nodes.Nodes[cfg.InstanceID]

	insClient := &tencent.InstanceTencent{}
	res, err := insClient.CreateInstance(1)
	if err != nil {
		return nil, err
	}
	nodes.AddNodeFromInstanceIdSet(res, clusterName)

	return n.Delete(ctx, cfg)
}
