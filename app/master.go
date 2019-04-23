package app

import (
	"context"
	pb "mobingi/ocean/app/proto/master"
)

type master struct{}

func (m *master) Join(ctx context.Context, cfg *pb.MasterConfig) (*pb.JoinResponse, error) {
	return &pb.JoinResponse{State: "success", Message: ""}, nil
}
