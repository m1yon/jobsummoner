package http

import (
	"context"
	"io"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/a-h/templ"
	"github.com/m1yon/jobsummoner"
	"github.com/m1yon/jobsummoner/internal/components"
	"github.com/m1yon/jobsummoner/internal/job"
)

type Server struct {
	logger     *slog.Logger
	Render     func(component templ.Component, ctx context.Context, w io.Writer) error
	jobService jobsummoner.JobService
	*http.Server
}

func NewServer(logger *slog.Logger, jobService *job.DefaultJobService) *Server {
	s := &Server{
		logger:     logger,
		Render:     components.Render,
		jobService: jobService,
	}

	s.Server = &http.Server{
		Addr:         ":3000",
		Handler:      s.routes(),
		ErrorLog:     slog.NewLogLogger(s.logger.Handler(), slog.LevelError),
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	return s
}

func (s *Server) Start(addr string) {
	s.logger.Info("server started", "addr", s.Server.Addr)

	err := s.ListenAndServe()
	s.logger.Error(err.Error())
	os.Exit(1)
}
