package interactions

import "fmt"

func Curse(name string) string {
	return fmt.Sprintf("Go to hell, %s!", name)
}

type CurseAdapter func(name string) string

func (g CurseAdapter) Curse(name string) (string, error) {
	return g(name), nil
}
