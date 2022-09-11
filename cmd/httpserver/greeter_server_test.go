package main_test

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/docker/go-connections/nat"
	"github.com/quii/go-specs-greet/adapters"
	go_specs_greet "github.com/quii/go-specs-greet/adapters/httpserver"
	"github.com/quii/go-specs-greet/specifications"
	"github.com/testcontainers/testcontainers-go/wait"
)

func TestGreeterServer(t *testing.T) {
	ctx := context.Background()
	port := "8080"

	adapters.StartDockerServer(
		t,
		ctx,
		"./cmd/httpserver/Dockerfile",
		port,
		wait.ForListeningPort(nat.Port(port)).WithStartupTimeout(5*time.Second),
	)

	client := http.Client{
		Timeout: 1 * time.Second,
	}

	driver := go_specs_greet.Driver{BaseURL: "http://localhost:8080", Client: &client}
	specifications.GreetSpecification(t, driver)
}
