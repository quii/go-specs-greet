package main

import (
	"log"
	"net/http"

	"github.com/quii/go-specs-greet/adapters/webserver"
)

func main() {
	handler, err := webserver.NewHandler()
	if err != nil {
		log.Fatal(err)
	}
	log.Fatal(http.ListenAndServe(":8081", handler))
}
