package grpcserver

import (
	"context"

	go_specs_greet "github.com/quii/go-specs-greet"
)

type GreetServer struct {
	UnimplementedGreeterServer
}

func (g GreetServer) Greet(ctx context.Context, request *GreetRequest) (*GreetReply, error) {
	return &GreetReply{Message: go_specs_greet.Greet(request.Name)}, nil
}
