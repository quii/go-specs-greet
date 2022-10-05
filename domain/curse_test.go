package domain_test

import (
	"testing"

	"github.com/quii/go-specs-greet/domain"
	"github.com/quii/go-specs-greet/specifications"
)

func TestCurse(t *testing.T) {
	specifications.CurseSpecification(
		t,
		domain.CurseAdapter(domain.Curse),
	)
}
