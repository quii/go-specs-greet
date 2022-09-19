package httpserver

import (
	"fmt"
	"net/http"

	"github.com/quii/go-specs-greet"
)

const (
	greetPath = "/greet"
	cursePath = "/curse"
)

func NewHandler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc(greetPath, replyWith(gospecsgreet.Greet))
	mux.HandleFunc(cursePath, replyWith(gospecsgreet.Curse))
	return mux
}

func replyWith(f func(name string) (interaction string)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		name := r.URL.Query().Get("name")
		fmt.Fprint(w, f(name))
	}
}
