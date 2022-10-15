package pages

import (
	"fmt"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/input"
)

type Form struct {
	Page *rod.Page
}

func (f Form) Greet(name string) error {
	greetInput, err := f.Page.Element("#greet-input")
	if err != nil {
		return fmt.Errorf("couldn't find #greet-input on Page")
	}
	return greetInput.MustInput(name).Type(input.Enter)
}

func (f Form) Curse(name string) error {
	curseInput, err := f.Page.Element("#curse-input")
	if err != nil {
		return fmt.Errorf("couldn't find #curse-input on Page")
	}
	return curseInput.MustInput(name).Type(input.Enter)
}
