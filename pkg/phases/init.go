package phases

import (
	"mobingi/ocean/pkg/config"
	configstorage "mobingi/ocean/pkg/storage"
)

func Init(cfg *config.Config) (configstorage.Cluster, error) {
	storage, err := configstorage.NewStorage(&configstorage.ClusterMongo{}, cfg)
	if err != nil {
		return nil, err
	}
	return storage, nil
}
