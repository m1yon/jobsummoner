package main

import (
	"log"
	"net/http"

	"github.com/m1yon/jobsummoner"
)

func main() {
	server := jobsummoner.DefaultServer{}
	handler := http.HandlerFunc(server.ServerHTTP)
	log.Fatal(http.ListenAndServe(":3000", handler))
}
