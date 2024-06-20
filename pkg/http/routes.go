package http

import "net/http"

func (server *DefaultServer) routes() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /{$}", server.getHomepageHandler)

	return mux
}
