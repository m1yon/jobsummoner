package main

import (
	"log"
	"net/http"
)

func main() {
	handler := http.HandlerFunc(HomepageServer)
	log.Fatal(http.ListenAndServe(":3000", handler))
}
