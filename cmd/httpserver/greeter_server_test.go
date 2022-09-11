package main_test

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/quii/go-specs-greet/adapters"
	go_specs_greet "github.com/quii/go-specs-greet/adapters/httpserver"
	"github.com/quii/go-specs-greet/specifications"
)

func TestGreeterServer(t *testing.T) {
	var (
		ctx            = context.Background()
		port           = "8080"
		dockerFilePath = "./cmd/httpserver/Dockerfile"
		baseURL        = fmt.Sprintf("http://localhost:%s", port)
		driver         = go_specs_greet.Driver{BaseURL: baseURL, Client: &http.Client{
			Timeout: 1 * time.Second,
		}}
	)

	adapters.StartDockerServer(ctx, t, dockerFilePath, port)
	specifications.GreetSpecification(t, driver)
}
