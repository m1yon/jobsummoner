package main

import "net/http"

func main() {
	handler := http.HandlerFunc(HomepageServer)
	http.ListenAndServe(":3000", handler)
}
