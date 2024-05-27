package main

import "net/http"

func HomepageServer(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello, world!"))
}
