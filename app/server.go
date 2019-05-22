package app

import (
	pb "mobingi/ocean/app/proto"
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

	pb.RegisterMasterServer(s, &master{})
	pb.RegisterClusterServer(s, &cluster{})
	pb.RegisterNodeServer(s, &node{})
	if err := s.Serve(lis); err != nil {
		log.Errorf("failed to serve: %v", err)
	}
	return nil
}
