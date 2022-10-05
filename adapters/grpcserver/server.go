package grpcserver

import (
	"context"

	"github.com/quii/go-specs-greet/domain"
)

type GreetServer struct {
	UnimplementedGreeterServer
}

func (g GreetServer) Curse(ctx context.Context, request *CurseRequest) (*CurseReply, error) {
	return &CurseReply{Message: domain.Curse(request.Name)}, nil
}

func (g GreetServer) Greet(ctx context.Context, request *GreetRequest) (*GreetReply, error) {
	return &GreetReply{Message: domain.Greet(request.Name)}, nil
}
