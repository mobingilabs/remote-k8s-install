package client

import (
	"mobingi/ocean/pkg/storage"
)

var Clusters []string

func InitClustersAndNodes() error {
	storage := storage.NewStorage()
	clusterArr, err := storage.All()
	if err != nil {
		return err
	}
	for _, name := range clusterArr {
		Clusters = append(Clusters, name)

		node, err := NewNodeClient(name)
		if err != nil {
			return err
		}
		nodeList, err := node.GetNode()
		for _, n := range nodeList.Items {
			Nodes[n.GetName()] = name
		}
	}
	return nil
}

// TODO 自动检测新集群加入监视
func ClustersMonitor() error {
	for _, cluster := range Clusters {
		node, err := NewNodeClient(cluster)
		if err != nil {
			return err
		}
		node.NewUnhealthyNodeTimer()
	}
	return nil
}
