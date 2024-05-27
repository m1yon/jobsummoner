package main

import (
	"log"
	"net/http"
)

func main() {
	server := DefaultHomepageServer{}
	handler := http.HandlerFunc(server.ServerHTTP)
	log.Fatal(http.ListenAndServe(":3000", handler))
}
