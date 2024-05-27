package main

import (
	"log"
	"net/http"
)

func main() {
	server := DefaultHomepageServer{}
	handler := http.HandlerFunc(server.Get)
	log.Fatal(http.ListenAndServe(":3000", handler))
}
