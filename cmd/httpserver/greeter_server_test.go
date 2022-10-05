package main_test

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/quii/go-specs-greet/adapters"
	go_specs_greet "github.com/quii/go-specs-greet/adapters/httpserver"
	"github.com/quii/go-specs-greet/specifications"
)

func TestGreeterServer(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}
	var (
		port   = "8080"
		driver = go_specs_greet.TestDriver{
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
