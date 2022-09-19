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
	if testing.Short() {
		t.Skip()
	}
	var (
		ctx    = context.Background()
		port   = "50051"
		addr   = fmt.Sprintf("localhost:%s", port)
		driver = grpcserver.Driver{Addr: addr}
	)

	t.Cleanup(driver.Close)
	adapters.StartDockerServer(ctx, t, port, "grpcserver")
	specifications.GreetSpecification(t, &driver)
	specifications.CurseSpecification(t, &driver)
}
