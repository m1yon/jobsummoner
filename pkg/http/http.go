package http

import (
	"context"
	"database/sql"
	"io"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/a-h/templ"
	"github.com/alexedwards/scs/sqlite3store"
	"github.com/alexedwards/scs/v2"
	"github.com/go-playground/form/v4"
	"github.com/m1yon/jobsummoner"
	"github.com/m1yon/jobsummoner/internal/components"
)

type Server struct {
	logger         *slog.Logger
	Render         func(component templ.Component, ctx context.Context, w io.Writer) error
	jobService     jobsummoner.JobService
	userService    jobsummoner.UserService
	sessionManager *scs.SessionManager
	formDecoder    *form.Decoder
	*http.Server
}

func NewServer(logger *slog.Logger, jobService jobsummoner.JobService, userService jobsummoner.UserService, db *sql.DB) *Server {
	formDecoder := form.NewDecoder()

	sessionManager := scs.New()
	sessionManager.Store = sqlite3store.New(db)
	sessionManager.Lifetime = 12 * time.Hour
	sessionManager.Cookie.Secure = os.Getenv("FLY_APP_NAME") != ""

	s := &Server{
		logger:         logger,
		Render:         components.Render,
		jobService:     jobService,
		userService:    userService,
		sessionManager: sessionManager,
		formDecoder:    formDecoder,
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
