package main_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/quii/go-specs-greet/adapters"
	"github.com/quii/go-specs-greet/adapters/grpcserver"
	"github.com/quii/go-specs-greet/specifications"
)

func TestGreeterServer(t *testing.T) {
	var (
		ctx            = context.Background()
		port           = "50051"
		dockerFilePath = "./cmd/grpcserver/Dockerfile"
		addr           = fmt.Sprintf("localhost:%s", port)
		driver         = grpcserver.Driver{Addr: addr}
	)

	adapters.StartDockerServer(ctx, t, dockerFilePath, port)
	specifications.GreetSpecification(t, &driver)
}
