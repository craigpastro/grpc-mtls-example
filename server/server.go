package server

import (
	"context"

	pb "github.com/craigpastro/mtls-example/protos/api/v1"
)

type server struct {
	pb.UnimplementedServiceServer
}

func NewServer() *server {
	return &server{}
}

func (s *server) Echo(ctx context.Context, in *pb.EchoRequest) (*pb.EchoResponse, error) {
	return &pb.EchoResponse{
		Said: in.Say,
	}, nil
}
