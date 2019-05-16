package client

import (
	"context"
	"errors"
	"mobingi/ocean/pkg/kubernetes/client/nodes"
	"mobingi/ocean/pkg/log"
	"mobingi/ocean/pkg/storage"
)

var Clusters []string
var Monitors = make(map[string]context.CancelFunc)

func InitClustersAndNodes() error {
	storage := storage.NewStorage()
	clusterArr, err := storage.All()
	if err != nil {
		return err
	}
	for _, name := range clusterArr {
		node, err := NewNodeClient(name)
		if err != nil {
			return err
		}
		nodeList, err := node.GetNode()
		if err != nil {
			log.Errorf("%s: "+err.Error(), name)
			continue
			// return err
		}
		for _, n := range nodeList.Items {
			nodes.Nodes[n.GetName()] = name
		}
		Clusters = append(Clusters, name)
	}
	return nil
}

func InitClustersMonitor() error {
	for _, cluster := range Clusters {
		CreateMonitor(cluster)
	}
	return nil
}

func CreateMonitor(cluster string) error {
	node, err := NewNodeClient(cluster)
	if err != nil {
		return err
	}
	ctx, cancel := context.WithCancel(context.Background())
	Monitors[cluster] = cancel
	node.NewUnhealthyNodeTimer(ctx)
	return nil
}

func StopMonitor(cluster string) error {
	cancel, ok := Monitors[cluster]
	if !ok {
		return errors.New("Not found")
	}
	cancel()
	return nil
}
