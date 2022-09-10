package go_specs_greet_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	go_specs_greet "github.com/quii/go-specs-greet"
	"github.com/quii/go-specs-greet/specifications"
)

func TestGreeterHandler(t *testing.T) {
	client := http.Client{
		Timeout: 5 * time.Second,
	}
	server := httptest.NewServer(http.HandlerFunc(go_specs_greet.Handler))
	defer server.Close()
	driver := go_specs_greet.Driver{BaseURL: server.URL, Client: &client}
	specifications.GreetSpecification(t, driver)
}
