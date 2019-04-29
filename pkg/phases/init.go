package phases

import (
	"mobingi/ocean/pkg/config"
	configstorage "mobingi/ocean/pkg/storage"
)

func Init(cfg *config.Config) error {
	_, err := configstorage.NewStorage(&configstorage.ClusterMongo{}, cfg)
	if err != nil {
		return err
	}
	return nil
}
