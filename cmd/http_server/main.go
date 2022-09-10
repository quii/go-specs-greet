package main

import (
	"net/http"

	go_specs_greet "github.com/quii/go-specs-greet/adapters/http"
)

func main() {
	handler := http.HandlerFunc(go_specs_greet.Handler)
	http.ListenAndServe(":8080", handler)
}
