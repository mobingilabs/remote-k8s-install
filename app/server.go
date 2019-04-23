package app

import (
	pbmaster "mobingi/ocean/app/proto/master"
	"mobingi/ocean/pkg/log"
	"net"

	"google.golang.org/grpc"
)

const (
	port = ":50051"
)

func Start() error {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Errorf("failed to listen: %v", err)
	}
	s := grpc.NewServer()

	pbmaster.RegisterMasterServer(s, &master{})
	if err := s.Serve(lis); err != nil {
		log.Errorf("failed to serve: %v", err)
	}

	// cfg, err := config.LoadConfigFromFile("config.yaml")
	// if err != nil {
	// 	return err
	// }
	// if err := master.InstallMasters(cfg); err != nil {
	// 	return err
	// }

	// adminConf, _ := cache.GetOne(constants.KubeconfPrefix, "admin.conf")

	// bootstrapconf, err := bootstrap.Bootstrap(adminConf.([]byte))
	// if err != nil {
	// 	log.Panic(err)
	// }

	// mi := &machine.MachineInfo{
	// 	PublicIP: cfg.Nodes[0].PublicIP,
	// 	User:     cfg.Nodes[0].User,
	// 	Password: cfg.Nodes[0].Password,
	// }

	// if err := node.Join(adminConf.([]byte), bootstrapconf, cfg.DownloadBinSite, mi); err != nil {
	// 	return err
	// }
	// mi.PublicIP = cfg.Nodes[1].PublicIP
	// mi.User = cfg.Nodes[1].User
	// mi.Password = cfg.Nodes[0].Password
	// if err := node.Join(adminConf.([]byte), bootstrapconf, cfg.DownloadBinSite, mi); err != nil {
	// 	return err
	// }

	return nil
}
