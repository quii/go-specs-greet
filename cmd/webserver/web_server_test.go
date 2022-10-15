package main_test

import (
	"fmt"
	"testing"

	"github.com/alecthomas/assert/v2"
	"github.com/quii/go-specs-greet/adapters"
	"github.com/quii/go-specs-greet/adapters/webserver"
	"github.com/quii/go-specs-greet/specifications"
)

func TestGreeterWeb(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}
	var (
		port            = "8081"
		driver, cleanup = webserver.NewDriver(fmt.Sprintf("http://localhost:%s", port))
	)

	t.Cleanup(func() {
		assert.NoError(t, cleanup())
	})

	adapters.StartDockerServer(t, port, "webserver")
	specifications.GreetSpecification(t, driver)
	specifications.CurseSpecification(t, driver)
}
