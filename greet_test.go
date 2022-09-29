package gospecsgreet_test

import (
	"testing"

	"github.com/alecthomas/assert/v2"
	go_specs_greet "github.com/quii/go-specs-greet"
	"github.com/quii/go-specs-greet/specifications"
)

func TestGreet(t *testing.T) {
	specifications.GreetSpecification(
		t,
		go_specs_greet.GreetAdapter(go_specs_greet.Greet),
	)

	t.Run("default name to world if it's an empty string", func(t *testing.T) {
		assert.Equal(t, "Hello, World", go_specs_greet.Greet(""))
	})
}
