package main

import (
	"github.com/m1yon/jobsummoner/pkg/http"
)

func main() {
	server := http.NewDefaultServer()
	server.ListenAndServe(":3000")
}
