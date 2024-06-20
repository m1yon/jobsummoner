package http

import (
	"context"
	"io"
	"log"
	"log/slog"
	"net/http"

	"github.com/a-h/templ"
	"github.com/m1yon/jobsummoner"
	"github.com/m1yon/jobsummoner/internal/components"
	"github.com/m1yon/jobsummoner/internal/job"
)

type Server interface {
	Render(component templ.Component, ctx context.Context, w io.Writer) error
	ListenAndServe(addr string)
}

type DefaultServer struct {
	logger *slog.Logger
	Render func(component templ.Component, ctx context.Context, w io.Writer) error
	jobsummoner.JobService
}

func NewDefaultServer(logger *slog.Logger, jobService *job.DefaultJobService) *DefaultServer {
	return &DefaultServer{
		logger:     logger,
		Render:     components.Render,
		JobService: jobService,
	}
}

func (server *DefaultServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	server.routes().ServeHTTP(w, r)
}

func (server *DefaultServer) ListenAndServe(addr string) {
	log.Fatal(http.ListenAndServe(addr, server))
}
