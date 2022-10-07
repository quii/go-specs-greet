package specifications

type CurseAdapter func(name string) string

func (g CurseAdapter) Curse(name string) (string, error) {
	return g(name), nil
}

type GreetAdapter func(name string) string

func (g GreetAdapter) Greet(name string) (string, error) {
	return g(name), nil
}
