package jobsummoner

type ScrapedJobsResults struct {
	Jobs []Job
}

type Scraper interface {
	ScrapeJobs() (ScrapedJobsResults, []error)
}

type ScrapeService interface {
	Start(scraper Scraper, crontab string)
}
