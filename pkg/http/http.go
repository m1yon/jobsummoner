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

type Server struct {
	logger     *slog.Logger
	Render     func(component templ.Component, ctx context.Context, w io.Writer) error
	jobService jobsummoner.JobService
}

func NewServer(logger *slog.Logger, jobService *job.DefaultJobService) *Server {
	return &Server{
		logger:     logger,
		Render:     components.Render,
		jobService: jobService,
	}
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.routes().ServeHTTP(w, r)
}

func (s *Server) Start(addr string) {
	s.logger.Info("server started", "port", "3000")
	log.Fatal(http.ListenAndServe(addr, s))
}
