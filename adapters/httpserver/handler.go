package httpserver

import (
	"fmt"
	"net/http"

	"github.com/quii/go-specs-greet"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	fmt.Fprint(w, gospecsgreet.Greet(name))
}
