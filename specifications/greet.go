package specifications

import (
	"testing"

	"github.com/alecthomas/assert/v2"
)

type Greeter interface {
	Greet(name string) (string, error)
}

func GreetSpecification(t *testing.T, greeter Greeter) {
	t.Run("greets a person", func(t *testing.T) {
		got, err := greeter.Greet("Mike")
		assert.NoError(t, err)
		assert.Equal(t, got, "Hello, Mike")
	})

	t.Run("when no name is supplied, greet the world", func(t *testing.T) {
		got, err := greeter.Greet("")
		assert.NoError(t, err)
		assert.Equal(t, got, "Hello, World")
	})
}
