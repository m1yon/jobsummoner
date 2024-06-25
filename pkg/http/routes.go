package http

import (
	"net/http"

	"github.com/justinas/alice"
)

func (s *Server) routes() http.Handler {
	mux := http.NewServeMux()

	dynamic := alice.New(s.sessionManager.LoadAndSave)

	mux.Handle("GET /{$}", dynamic.ThenFunc(s.getHomepageHandler))

	standard := alice.New(s.logRequest, s.recoverPanic, commonHeaders)
	return standard.Then(mux)
}
