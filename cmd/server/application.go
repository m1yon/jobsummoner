package main

import (
	"database/sql"
	"log/slog"
	"os"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/go-playground/form/v4"
	"github.com/m1yon/jobsummoner/internal/models"
	"github.com/m1yon/jobsummoner/internal/sqlite3store"
)

type application struct {
	logger         *slog.Logger
	jobs           models.JobModelInterface
	users          models.UserModelInterface
	sessionManager *scs.SessionManager
	formDecoder    *form.Decoder
}

func newApplication(logger *slog.Logger, jobs models.JobModelInterface, users models.UserModelInterface, db *sql.DB) *application {
	formDecoder := form.NewDecoder()

	sessionManager := scs.New()
	sessionManager.Store = sqlite3store.New(db)
	sessionManager.Lifetime = 12 * time.Hour
	sessionManager.Cookie.Secure = os.Getenv("FLY_APP_NAME") != ""

	app := &application{
		logger:         logger,
		jobs:           jobs,
		users:          users,
		sessionManager: sessionManager,
		formDecoder:    formDecoder,
	}

	return app
}
