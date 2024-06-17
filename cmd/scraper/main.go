package main

import (
	"database/sql"
	"log/slog"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/jonboulle/clockwork"
	"github.com/lmittmann/tint"
	"github.com/m1yon/jobsummoner"
	"github.com/m1yon/jobsummoner/internal/company"
	"github.com/m1yon/jobsummoner/internal/job"
	"github.com/m1yon/jobsummoner/internal/scrape"
	"github.com/m1yon/jobsummoner/internal/sqlitedb"
	"github.com/m1yon/jobsummoner/pkg/linkedin"

	_ "github.com/tursodatabase/libsql-client-go/libsql"
	_ "modernc.org/sqlite"
)

func main() {
	logger := slog.New(tint.NewHandler(os.Stdout, &tint.Options{Level: slog.LevelDebug}))
	err := godotenv.Load()

	if err != nil {
		logger.Warn("no .env file found")
	}

	httpClient := initHttpClient(logger)

	db, err := openDB(logger)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	scrapeService := initScrapeService(logger, db)

	scrapers := []jobsummoner.Scraper{
		linkedin.NewLinkedInJobScraper(linkedin.NewHttpLinkedInReader(linkedin.LinkedInReaderConfig{
			Keywords: []string{"typescript"},
			Location: "United States",
		}, httpClient, logger), logger),
	}

	scrapeService.Start(scrapers, "TZ=America/Denver */30 7-22 * * *", true)
}

func initHttpClient(logger *slog.Logger) *http.Client {
	proxyConfig := linkedin.ProxyConfig{
		Hostname: os.Getenv("PROXY_HOSTNAME"),
		Port:     os.Getenv("PROXY_PORT"),
		Username: os.Getenv("PROXY_USERNAME"),
		Password: os.Getenv("PROXY_PASSWORD"),
	}

	httpClient, err := linkedin.NewHttpProxyClientFromConfig(proxyConfig)

	if err != nil {
		logger.Warn("proxy server disabled", "reason", err.Error())
	} else {
		logger.Info("proxy server enabled")
	}

	return httpClient
}

func openDB(logger *slog.Logger) (*sql.DB, error) {
	db, err := sqlitedb.NewDB(logger, &sqlitedb.SqlConnectionOpener{})

	if err != nil {
		return nil, err
	}

	err = db.Ping()

	if err != nil {
		return nil, err
	}

	return db, nil
}

func initScrapeService(logger *slog.Logger, db *sql.DB) *scrape.DefaultScrapeService {
	c := clockwork.NewRealClock()

	companyRepository := sqlitedb.NewSqliteCompanyRepository(db)
	companyService := company.NewDefaultCompanyService(companyRepository)
	jobRepository := sqlitedb.NewSqliteJobRepository(db)
	jobService := job.NewDefaultJobService(jobRepository, companyService)

	scrapeRepository := sqlitedb.NewSqliteScrapeRepository(db, c)
	scrapeService := scrape.NewDefaultScrapeService(c, logger, scrapeRepository, jobService)

	return scrapeService
}
