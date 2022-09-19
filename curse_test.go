package gospecsgreet_test

import (
	"testing"

	go_specs_greet "github.com/quii/go-specs-greet"
	"github.com/quii/go-specs-greet/specifications"
)

func TestCurse(t *testing.T) {
	specifications.CurseSpecification(
		t,
		go_specs_greet.CurseAdapter(go_specs_greet.Curse),
	)
}
