package main

import (
	"database/sql"
	"log/slog"
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

	scrapers := []jobsummoner.Scraper{
		linkedin.NewLinkedInJobScraper(linkedin.NewHttpLinkedInReader(linkedin.LinkedInReaderConfig{
			Keywords: []string{"typescript"},
			Location: "United States",
		}, httpClient, logger), logger),
	}

	c := clockwork.NewRealClock()

	db, err := sqlitedb.NewDB(sql.Open)

	if err != nil {
		logger.Error("failed starting db")
		os.Exit(1)
	}

	err = db.Ping()

	if err != nil {
		logger.Error("failed pinging db")
		os.Exit(1)
	}

	companyRepository := sqlitedb.NewSqliteCompanyRepository(db)
	companyService := company.NewDefaultCompanyService(companyRepository)
	jobRepository := sqlitedb.NewSqliteJobRepository(db)
	jobService := job.NewDefaultJobService(jobRepository, companyService)

	scrapeRepository := sqlitedb.NewSqliteScrapeRepository(db, c)
	scrapeService := scrape.NewDefaultScrapeService(c, logger, scrapeRepository, jobService)

	scrapeService.Start(scrapers, "TZ=America/Denver */30 7-22 * * *", true)
}
