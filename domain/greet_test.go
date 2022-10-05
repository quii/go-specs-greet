package domain_test

import (
	"testing"

	"github.com/alecthomas/assert/v2"
	"github.com/quii/go-specs-greet/domain"
	"github.com/quii/go-specs-greet/specifications"
)

type GreetAdapter func(name string) string

func (g GreetAdapter) Greet(name string) (string, error) {
	return g(name), nil
}


func TestGreet(t *testing.T) {
	specifications.GreetSpecification(
		t,
		GreetAdapter(domain.Greet),
	)

	t.Run("default name to world if it's an empty string", func(t *testing.T) {
		assert.Equal(t, "Hello, World", domain.Greet(""))
	})
}
