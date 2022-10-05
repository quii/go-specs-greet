package httpserver

import (
	"fmt"
	"net/http"

	"github.com/quii/go-specs-greet/domain"
)

const (
	greetPath = "/greet"
	cursePath = "/curse"
)

func NewHandler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc(greetPath, replyWith(domain.Greet))
	mux.HandleFunc(cursePath, replyWith(domain.Curse))
	return mux
}

func replyWith(f func(name string) (interaction string)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		name := r.URL.Query().Get("name")
		fmt.Fprint(w, f(name))
	}
}
