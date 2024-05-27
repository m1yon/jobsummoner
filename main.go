package main

import (
	"log"
	"net/http"
)

func main() {
	server := DefaultServer{}
	handler := http.HandlerFunc(server.ServerHTTP)
	log.Fatal(http.ListenAndServe(":3000", handler))
}
