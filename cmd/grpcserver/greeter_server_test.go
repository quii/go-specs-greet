package main_test

import (
	"context"
	"testing"
	"time"

	"github.com/docker/go-connections/nat"
	"github.com/quii/go-specs-greet/adapters"
	"github.com/quii/go-specs-greet/adapters/grpcserver"
	"github.com/quii/go-specs-greet/specifications"
	"github.com/testcontainers/testcontainers-go/wait"
)

func TestGreeterServer(t *testing.T) {
	ctx := context.Background()
	port := "50051"

	adapters.StartDockerServer(ctx, t, "./cmd/grpcserver/Dockerfile", port, wait.ForListeningPort(nat.Port(port)).WithStartupTimeout(5*time.Second))

	driver := grpcserver.Driver{Addr: "localhost:50051"}
	specifications.GreetSpecification(t, &driver)
}
