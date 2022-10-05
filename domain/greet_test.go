package domain_test

import (
	"testing"

	"github.com/alecthomas/assert/v2"
	"github.com/quii/go-specs-greet/domain"
	"github.com/quii/go-specs-greet/specifications"
)

func TestGreet(t *testing.T) {
	specifications.GreetSpecification(
		t,
		domain.GreetAdapter(domain.Greet),
	)

	t.Run("default name to world if it's an empty string", func(t *testing.T) {
		assert.Equal(t, "Hello, World", domain.Greet(""))
	})
}
