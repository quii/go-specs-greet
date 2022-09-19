package grpcserver

import (
	"context"

	gospecsgreet "github.com/quii/go-specs-greet"
)

type GreetServer struct {
	UnimplementedGreeterServer
}

func (g GreetServer) Curse(ctx context.Context, request *CurseRequest) (*CurseReply, error) {
	return &CurseReply{Message: gospecsgreet.Curse(request.Name)}, nil
}

func (g GreetServer) Greet(ctx context.Context, request *GreetRequest) (*GreetReply, error) {
	return &GreetReply{Message: gospecsgreet.Greet(request.Name)}, nil
}
