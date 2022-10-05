package main

import (
	"log"
	"net/http"

	adapter "github.com/quii/go-specs-greet/adapters/httpserver"
)

func main() {
	if err := http.ListenAndServe(":8080", adapter.NewHandler()); err != nil {
		log.Fatal(err)
	}
}
