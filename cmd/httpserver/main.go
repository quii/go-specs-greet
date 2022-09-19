package main

import (
	"log"
	"net/http"

	go_specs_greet "github.com/quii/go-specs-greet/adapters/httpserver"
)

func main() {
	if err := http.ListenAndServe(":8080", go_specs_greet.NewHandler()); err != nil {
		log.Fatal(err)
	}
}
