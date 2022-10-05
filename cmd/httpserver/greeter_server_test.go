package main_test

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/quii/go-specs-greet/adapters"
	adapter "github.com/quii/go-specs-greet/adapters/httpserver"
	"github.com/quii/go-specs-greet/specifications"
)

func TestGreeterServer(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}
	var (
		port   = "8080"
		driver = adapter.TestDriver{
			BaseURL: fmt.Sprintf("http://localhost:%s", port),
			Client: &http.Client{
				Timeout: 1 * time.Second,
			},
		}
	)

	adapters.StartDockerServer(t, port, "httpserver")
	specifications.GreetSpecification(t, driver)
	specifications.CurseSpecification(t, driver)
}
