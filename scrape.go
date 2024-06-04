package jobsummoner

type Scraper interface {
	ScrapeJobs() ([]Job, []error)
}

type ScrapeService interface {
	Start(scraper []Scraper, crontab string)
}
