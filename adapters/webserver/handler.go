package webserver

import (
	_ "embed"
	"fmt"
	"net/http"

	"github.com/quii/go-specs-greet/domain/interactions"
)

const (
	greetPath = "/greet"
	cursePath = "/curse"
)

//go:embed form.html
var formMarkup string

func NewHandler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/", form)
	mux.HandleFunc(greetPath, replyWith(interactions.Greet))
	mux.HandleFunc(cursePath, replyWith(interactions.Curse))
	return mux
}

func replyWith(interact func(name string) string) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		fmt.Fprintf(
			w,
			`<h1 id="reply">%s</h1>`,
			interact(r.Form.Get("name")),
		)
	}
}

func form(w http.ResponseWriter, _ *http.Request) {
	fmt.Fprint(w, formMarkup)
}
