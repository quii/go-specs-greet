package domain_test

import (
	"testing"

	"github.com/quii/go-specs-greet/domain"
	"github.com/quii/go-specs-greet/specifications"
)

type CurseAdapter func(name string) string

func (g CurseAdapter) Curse(name string) (string, error) {
	return g(name), nil
}

func TestCurse(t *testing.T) {
	specifications.CurseSpecification(
		t,
		CurseAdapter(domain.Curse),
	)
}
