package linkedin

import (
	"bytes"
	"fmt"
	"log/slog"
	"testing"
	"time"

	"github.com/m1yon/jobsummoner"
	"github.com/stretchr/testify/assert"
)

func TestLinkedInScraper(t *testing.T) {
	t.Run("file reader - scrape and paginates job listings correctly with 2 pages", func(t *testing.T) {
		fileReader := NewFileLinkedInReader("./test-helpers/li-job-listings-%v.html")
		logger := slog.New(slog.NewTextHandler(nil, nil))
		scraper := NewCustomLinkedInJobScraper(fileReader, logger)

		got, errs := scraper.ScrapeJobs()
		want := []jobsummoner.Job{
			{Position: "Senior Software Engineer", CompanyID: "goliath-partners-inc", CompanyName: "Goliath Partners", Location: "New York City Metropolitan Area", URL: "https://www.linkedin.com/jobs/view/senior-software-engineer-at-goliath-partners-3941197019?position=1&pageNum=0&refId=QN9Xe6fr3eBaRYE4EMYCIg%3D%3D&trackingId=TtFssv3yCaLnbvL4n7pGqw%3D%3D&trk=public_jobs_jserp-result_search-card", SourceID: "linkedin"},
			{Position: "Full Stack Engineer", CompanyID: "goliath-partners-inc", CompanyName: "Goliath Partners", Location: "New York City Metropolitan Area", URL: "https://www.linkedin.com/jobs/view/full-stack-engineer-at-goliath-partners-3941199093?position=2&pageNum=0&refId=QN9Xe6fr3eBaRYE4EMYCIg%3D%3D&trackingId=qbBBmbaU%2FGecHsilZU6rPg%3D%3D&trk=public_jobs_jserp-result_search-card", SourceID: "linkedin"},
			{Position: "Application Developer", CompanyID: "proliance-consulting", CompanyName: "Proliance Consulting", Location: "Tempe, AZ", URL: "https://www.linkedin.com/jobs/view/application-developer-at-proliance-consulting-3938872094?position=3&pageNum=0&refId=QN9Xe6fr3eBaRYE4EMYCIg%3D%3D&trackingId=yp1wpeQB6Yl5V0C%2FQp1Rsw%3D%3D&trk=public_jobs_jserp-result_search-card", SourceID: "linkedin"},
			{Position: "DevOps Engineer", CompanyID: "firstpro", CompanyName: "firstPRO, Inc", Location: "Portland, Maine Metropolitan Area", URL: "https://www.linkedin.com/jobs/view/devops-engineer-at-firstpro-inc-3941423455?position=4&pageNum=0&refId=QN9Xe6fr3eBaRYE4EMYCIg%3D%3D&trackingId=uFqmQan2fgvaePCw9yJi%2Bw%3D%3D&trk=public_jobs_jserp-result_search-card", SourceID: "linkedin"},
			{Position: "Software Developer 4", CompanyID: "oracle", CompanyName: "Oracle", Location: "Seattle, WA", URL: "https://www.linkedin.com/jobs/view/software-developer-4-at-oracle-3926234960?position=5&pageNum=0&refId=QN9Xe6fr3eBaRYE4EMYCIg%3D%3D&trackingId=V9okmTSo%2FMAajAtqw9viYQ%3D%3D&trk=public_jobs_jserp-result_search-card", SourceID: "linkedin"},
			{Position: "Senior Software Engineer - Full Stack", CompanyID: "titan-invest", CompanyName: "Titan", Location: "New York City Metropolitan Area", URL: "https://www.linkedin.com/jobs/view/senior-software-engineer-full-stack-at-titan-3904725366?position=6&pageNum=0&refId=QN9Xe6fr3eBaRYE4EMYCIg%3D%3D&trackingId=sL3X%2BQ7k7ebAi9hQ2MzRzw%3D%3D&trk=public_jobs_jserp-result_search-card", SourceID: "linkedin"},
			{Position: "Software Developer 4", CompanyID: "oracle", CompanyName: "Oracle", Location: "Austin, TX", URL: "https://www.linkedin.com/jobs/view/software-developer-4-at-oracle-3926236707?position=7&pageNum=0&refId=QN9Xe6fr3eBaRYE4EMYCIg%3D%3D&trackingId=DzcGOX9aApHdauzctJsjGg%3D%3D&trk=public_jobs_jserp-result_search-card", SourceID: "linkedin"},
			{Position: "Software Developer 4", CompanyID: "oracle", CompanyName: "Oracle", Location: "Redwood City, CA", URL: "https://www.linkedin.com/jobs/view/software-developer-4-at-oracle-3914083553?position=8&pageNum=0&refId=QN9Xe6fr3eBaRYE4EMYCIg%3D%3D&trackingId=grmTG%2FuW6DmSaBwEmGKlBQ%3D%3D&trk=public_jobs_jserp-result_search-card", SourceID: "linkedin"},
			{Position: "Senior Software Engineer", CompanyID: "rentspree", CompanyName: "RentSpree", Location: "Greater Seattle Area", URL: "https://www.linkedin.com/jobs/view/senior-software-engineer-at-rentspree-3938840276?position=9&pageNum=0&refId=QN9Xe6fr3eBaRYE4EMYCIg%3D%3D&trackingId=dkiwVGJIEqnYXvSqgvjO8w%3D%3D&trk=public_jobs_jserp-result_search-card", SourceID: "linkedin"},
			{Position: "Analyst, Developer Specialty 3", CompanyID: "tekwissen", CompanyName: "TekWissen ®", Location: "St Paul, MN", URL: "https://www.linkedin.com/jobs/view/analyst-developer-specialty-3-at-tekwissen-%C2%AE-3938822375?position=10&pageNum=0&refId=QN9Xe6fr3eBaRYE4EMYCIg%3D%3D&trackingId=cq1o%2Bo56V3Fdpe038vGwug%3D%3D&trk=public_jobs_jserp-result_search-card", SourceID: "linkedin"},
			{Position: "Senior Principal Engineer - Relocation to Australia", CompanyID: "rokt", CompanyName: "Rokt", Location: "San Francisco, CA", URL: "https://www.linkedin.com/jobs/view/senior-principal-engineer-relocation-to-australia-at-rokt-3905316034?position=1&pageNum=2&refId=O7f5RZQjEgD2Ze4PH0mAuQ%3D%3D&trackingId=b%2BQG0mQCaIcp0cHeoz5bTQ%3D%3D&trk=public_jobs_jserp-result_search-card", SourceID: "linkedin"},
			{Position: "Principal Software Engineer/Sr Principal Software Engineer", CompanyID: "oracle", CompanyName: "Oracle", Location: "Portland, OR", URL: "https://www.linkedin.com/jobs/view/principal-software-engineer-sr-principal-software-engineer-at-oracle-3911661146?position=2&pageNum=2&refId=O7f5RZQjEgD2Ze4PH0mAuQ%3D%3D&trackingId=bfNqZnkoWYbjj6bnZS47Tw%3D%3D&trk=public_jobs_jserp-result_search-card", SourceID: "linkedin"},
			{Position: "Principal Software Engineer/Sr Principal Software Engineer", CompanyID: "oracle", CompanyName: "Oracle", Location: "Los Angeles, CA", URL: "https://www.linkedin.com/jobs/view/principal-software-engineer-sr-principal-software-engineer-at-oracle-3911662017?position=3&pageNum=2&refId=O7f5RZQjEgD2Ze4PH0mAuQ%3D%3D&trackingId=Y%2F9rLjlZm8z8F0NfkPlUIg%3D%3D&trk=public_jobs_jserp-result_search-card", SourceID: "linkedin"},
			{Position: "Software Development Manager", CompanyID: "oracle", CompanyName: "Oracle", Location: "Seattle, WA", URL: "https://www.linkedin.com/jobs/view/software-development-manager-at-oracle-3900536996?position=4&pageNum=2&refId=O7f5RZQjEgD2Ze4PH0mAuQ%3D%3D&trackingId=kjOrd%2F6ZadeHijvsY1vfIg%3D%3D&trk=public_jobs_jserp-result_search-card", SourceID: "linkedin"},
			{Position: "Software Development Senior Manager", CompanyID: "oracle", CompanyName: "Oracle", Location: "Pleasanton, CA", URL: "https://www.linkedin.com/jobs/view/software-development-senior-manager-at-oracle-3919395308?position=5&pageNum=2&refId=O7f5RZQjEgD2Ze4PH0mAuQ%3D%3D&trackingId=CDB4Sem7WIv%2BkGAF4DPkHQ%3D%3D&trk=public_jobs_jserp-result_search-card", SourceID: "linkedin"},
			{Position: "Principal Software Engineer/Sr Principal Software Engineer", CompanyID: "oracle", CompanyName: "Oracle", Location: "Austin, TX", URL: "https://www.linkedin.com/jobs/view/principal-software-engineer-sr-principal-software-engineer-at-oracle-3911660279?position=6&pageNum=2&refId=O7f5RZQjEgD2Ze4PH0mAuQ%3D%3D&trackingId=F7zX0LB%2ByVaZ9IufjIAFiQ%3D%3D&trk=public_jobs_jserp-result_search-card", SourceID: "linkedin"},
			{Position: "Senior Software Developer - Data Center Infrastructure Team", CompanyID: "oracle", CompanyName: "Oracle", Location: "Seattle, WA", URL: "https://www.linkedin.com/jobs/view/senior-software-developer-data-center-infrastructure-team-at-oracle-3900542100?position=7&pageNum=2&refId=O7f5RZQjEgD2Ze4PH0mAuQ%3D%3D&trackingId=fWzqEG9OjDymYMdpTrT5zA%3D%3D&trk=public_jobs_jserp-result_search-card", SourceID: "linkedin"},
			{Position: "Software Development Senior Manager", CompanyID: "oracle", CompanyName: "Oracle", Location: "Redwood City, CA", URL: "https://www.linkedin.com/jobs/view/software-development-senior-manager-at-oracle-3921687498?position=8&pageNum=2&refId=O7f5RZQjEgD2Ze4PH0mAuQ%3D%3D&trackingId=hZQ%2B%2B3rtpKibl1ZmT9V4zQ%3D%3D&trk=public_jobs_jserp-result_search-card", SourceID: "linkedin"},
			{Position: "Staff Machine Learning Engineer - Seattle", CompanyID: "rokt", CompanyName: "Rokt", Location: "Seattle, WA", URL: "https://www.linkedin.com/jobs/view/staff-machine-learning-engineer-seattle-at-rokt-3905393698?position=9&pageNum=2&refId=O7f5RZQjEgD2Ze4PH0mAuQ%3D%3D&trackingId=792eL8NsCwvBl2h%2FjoxGQA%3D%3D&trk=public_jobs_jserp-result_search-card", SourceID: "linkedin"},
		}

		assert.Equal(t, want, got)
		assert.Equal(t, 0, len(errs))
	})

	t.Run("http reader - scrape and paginates job listings correctly with 2 pages", func(t *testing.T) {
		stubClient := NewStubClient()
		httpReader := NewHttpLinkedInReader(LinkedInReaderConfig{
			Keywords: []string{"Software Engineer", "Manager"},
			Location: "United States",
			MaxAge:   time.Hour * 4,
		}, stubClient)
		logger := slog.New(slog.NewTextHandler(nil, nil))
		scraper := NewCustomLinkedInJobScraper(httpReader, logger)

		got, errs := scraper.ScrapeJobs()
		want := []jobsummoner.Job{
			{Position: "Senior Software Engineer", CompanyID: "goliath-partners-inc", CompanyName: "Goliath Partners", Location: "New York City Metropolitan Area", URL: "https://www.linkedin.com/jobs/view/senior-software-engineer-at-goliath-partners-3941197019?position=1&pageNum=0&refId=QN9Xe6fr3eBaRYE4EMYCIg%3D%3D&trackingId=TtFssv3yCaLnbvL4n7pGqw%3D%3D&trk=public_jobs_jserp-result_search-card", SourceID: "linkedin"},
			{Position: "Full Stack Engineer", CompanyID: "goliath-partners-inc", CompanyName: "Goliath Partners", Location: "New York City Metropolitan Area", URL: "https://www.linkedin.com/jobs/view/full-stack-engineer-at-goliath-partners-3941199093?position=2&pageNum=0&refId=QN9Xe6fr3eBaRYE4EMYCIg%3D%3D&trackingId=qbBBmbaU%2FGecHsilZU6rPg%3D%3D&trk=public_jobs_jserp-result_search-card", SourceID: "linkedin"},
			{Position: "Application Developer", CompanyID: "proliance-consulting", CompanyName: "Proliance Consulting", Location: "Tempe, AZ", URL: "https://www.linkedin.com/jobs/view/application-developer-at-proliance-consulting-3938872094?position=3&pageNum=0&refId=QN9Xe6fr3eBaRYE4EMYCIg%3D%3D&trackingId=yp1wpeQB6Yl5V0C%2FQp1Rsw%3D%3D&trk=public_jobs_jserp-result_search-card", SourceID: "linkedin"},
			{Position: "DevOps Engineer", CompanyID: "firstpro", CompanyName: "firstPRO, Inc", Location: "Portland, Maine Metropolitan Area", URL: "https://www.linkedin.com/jobs/view/devops-engineer-at-firstpro-inc-3941423455?position=4&pageNum=0&refId=QN9Xe6fr3eBaRYE4EMYCIg%3D%3D&trackingId=uFqmQan2fgvaePCw9yJi%2Bw%3D%3D&trk=public_jobs_jserp-result_search-card", SourceID: "linkedin"},
			{Position: "Software Developer 4", CompanyID: "oracle", CompanyName: "Oracle", Location: "Seattle, WA", URL: "https://www.linkedin.com/jobs/view/software-developer-4-at-oracle-3926234960?position=5&pageNum=0&refId=QN9Xe6fr3eBaRYE4EMYCIg%3D%3D&trackingId=V9okmTSo%2FMAajAtqw9viYQ%3D%3D&trk=public_jobs_jserp-result_search-card", SourceID: "linkedin"},
			{Position: "Senior Software Engineer - Full Stack", CompanyID: "titan-invest", CompanyName: "Titan", Location: "New York City Metropolitan Area", URL: "https://www.linkedin.com/jobs/view/senior-software-engineer-full-stack-at-titan-3904725366?position=6&pageNum=0&refId=QN9Xe6fr3eBaRYE4EMYCIg%3D%3D&trackingId=sL3X%2BQ7k7ebAi9hQ2MzRzw%3D%3D&trk=public_jobs_jserp-result_search-card", SourceID: "linkedin"},
			{Position: "Software Developer 4", CompanyID: "oracle", CompanyName: "Oracle", Location: "Austin, TX", URL: "https://www.linkedin.com/jobs/view/software-developer-4-at-oracle-3926236707?position=7&pageNum=0&refId=QN9Xe6fr3eBaRYE4EMYCIg%3D%3D&trackingId=DzcGOX9aApHdauzctJsjGg%3D%3D&trk=public_jobs_jserp-result_search-card", SourceID: "linkedin"},
			{Position: "Software Developer 4", CompanyID: "oracle", CompanyName: "Oracle", Location: "Redwood City, CA", URL: "https://www.linkedin.com/jobs/view/software-developer-4-at-oracle-3914083553?position=8&pageNum=0&refId=QN9Xe6fr3eBaRYE4EMYCIg%3D%3D&trackingId=grmTG%2FuW6DmSaBwEmGKlBQ%3D%3D&trk=public_jobs_jserp-result_search-card", SourceID: "linkedin"},
			{Position: "Senior Software Engineer", CompanyID: "rentspree", CompanyName: "RentSpree", Location: "Greater Seattle Area", URL: "https://www.linkedin.com/jobs/view/senior-software-engineer-at-rentspree-3938840276?position=9&pageNum=0&refId=QN9Xe6fr3eBaRYE4EMYCIg%3D%3D&trackingId=dkiwVGJIEqnYXvSqgvjO8w%3D%3D&trk=public_jobs_jserp-result_search-card", SourceID: "linkedin"},
			{Position: "Analyst, Developer Specialty 3", CompanyID: "tekwissen", CompanyName: "TekWissen ®", Location: "St Paul, MN", URL: "https://www.linkedin.com/jobs/view/analyst-developer-specialty-3-at-tekwissen-%C2%AE-3938822375?position=10&pageNum=0&refId=QN9Xe6fr3eBaRYE4EMYCIg%3D%3D&trackingId=cq1o%2Bo56V3Fdpe038vGwug%3D%3D&trk=public_jobs_jserp-result_search-card", SourceID: "linkedin"},
			{Position: "Senior Principal Engineer - Relocation to Australia", CompanyID: "rokt", CompanyName: "Rokt", Location: "San Francisco, CA", URL: "https://www.linkedin.com/jobs/view/senior-principal-engineer-relocation-to-australia-at-rokt-3905316034?position=1&pageNum=2&refId=O7f5RZQjEgD2Ze4PH0mAuQ%3D%3D&trackingId=b%2BQG0mQCaIcp0cHeoz5bTQ%3D%3D&trk=public_jobs_jserp-result_search-card", SourceID: "linkedin"},
			{Position: "Principal Software Engineer/Sr Principal Software Engineer", CompanyID: "oracle", CompanyName: "Oracle", Location: "Portland, OR", URL: "https://www.linkedin.com/jobs/view/principal-software-engineer-sr-principal-software-engineer-at-oracle-3911661146?position=2&pageNum=2&refId=O7f5RZQjEgD2Ze4PH0mAuQ%3D%3D&trackingId=bfNqZnkoWYbjj6bnZS47Tw%3D%3D&trk=public_jobs_jserp-result_search-card", SourceID: "linkedin"},
			{Position: "Principal Software Engineer/Sr Principal Software Engineer", CompanyID: "oracle", CompanyName: "Oracle", Location: "Los Angeles, CA", URL: "https://www.linkedin.com/jobs/view/principal-software-engineer-sr-principal-software-engineer-at-oracle-3911662017?position=3&pageNum=2&refId=O7f5RZQjEgD2Ze4PH0mAuQ%3D%3D&trackingId=Y%2F9rLjlZm8z8F0NfkPlUIg%3D%3D&trk=public_jobs_jserp-result_search-card", SourceID: "linkedin"},
			{Position: "Software Development Manager", CompanyID: "oracle", CompanyName: "Oracle", Location: "Seattle, WA", URL: "https://www.linkedin.com/jobs/view/software-development-manager-at-oracle-3900536996?position=4&pageNum=2&refId=O7f5RZQjEgD2Ze4PH0mAuQ%3D%3D&trackingId=kjOrd%2F6ZadeHijvsY1vfIg%3D%3D&trk=public_jobs_jserp-result_search-card", SourceID: "linkedin"},
			{Position: "Software Development Senior Manager", CompanyID: "oracle", CompanyName: "Oracle", Location: "Pleasanton, CA", URL: "https://www.linkedin.com/jobs/view/software-development-senior-manager-at-oracle-3919395308?position=5&pageNum=2&refId=O7f5RZQjEgD2Ze4PH0mAuQ%3D%3D&trackingId=CDB4Sem7WIv%2BkGAF4DPkHQ%3D%3D&trk=public_jobs_jserp-result_search-card", SourceID: "linkedin"},
			{Position: "Principal Software Engineer/Sr Principal Software Engineer", CompanyID: "oracle", CompanyName: "Oracle", Location: "Austin, TX", URL: "https://www.linkedin.com/jobs/view/principal-software-engineer-sr-principal-software-engineer-at-oracle-3911660279?position=6&pageNum=2&refId=O7f5RZQjEgD2Ze4PH0mAuQ%3D%3D&trackingId=F7zX0LB%2ByVaZ9IufjIAFiQ%3D%3D&trk=public_jobs_jserp-result_search-card", SourceID: "linkedin"},
			{Position: "Senior Software Developer - Data Center Infrastructure Team", CompanyID: "oracle", CompanyName: "Oracle", Location: "Seattle, WA", URL: "https://www.linkedin.com/jobs/view/senior-software-developer-data-center-infrastructure-team-at-oracle-3900542100?position=7&pageNum=2&refId=O7f5RZQjEgD2Ze4PH0mAuQ%3D%3D&trackingId=fWzqEG9OjDymYMdpTrT5zA%3D%3D&trk=public_jobs_jserp-result_search-card", SourceID: "linkedin"},
			{Position: "Software Development Senior Manager", CompanyID: "oracle", CompanyName: "Oracle", Location: "Redwood City, CA", URL: "https://www.linkedin.com/jobs/view/software-development-senior-manager-at-oracle-3921687498?position=8&pageNum=2&refId=O7f5RZQjEgD2Ze4PH0mAuQ%3D%3D&trackingId=hZQ%2B%2B3rtpKibl1ZmT9V4zQ%3D%3D&trk=public_jobs_jserp-result_search-card", SourceID: "linkedin"},
			{Position: "Staff Machine Learning Engineer - Seattle", CompanyID: "rokt", CompanyName: "Rokt", Location: "Seattle, WA", URL: "https://www.linkedin.com/jobs/view/staff-machine-learning-engineer-seattle-at-rokt-3905393698?position=9&pageNum=2&refId=O7f5RZQjEgD2Ze4PH0mAuQ%3D%3D&trackingId=792eL8NsCwvBl2h%2FjoxGQA%3D%3D&trk=public_jobs_jserp-result_search-card", SourceID: "linkedin"},
		}

		assert.Equal(t, want, got)
		assert.Equal(t, 0, len(errs))
	})

	t.Run("handles invalid company URLs", func(t *testing.T) {
		mockReader := NewFileLinkedInReader("./test-helpers/li-job-listings_bad-company-url-%v.html")
		logBufferSpy := new(bytes.Buffer)
		logger := slog.New(slog.NewTextHandler(logBufferSpy, nil))
		scraper := NewCustomLinkedInJobScraper(mockReader, logger)

		got, errs := scraper.ScrapeJobs()
		want := []jobsummoner.Job{
			{Position: "Software Engineer II (Frontend) - Seller Experience", CompanyID: "stubhub", CompanyName: "StubHub", Location: "Los Angeles, CA", URL: "https://www.linkedin.com/jobs/view/software-engineer-ii-frontend-seller-experience-at-stubhub-3916280897?position=2&pageNum=0&refId=fsDMYm%2BoJB2zdtWm%2FhnZ3g%3D%3D&trackingId=dbosG%2Ftu2ZxD8zDrXnrTWw%3D%3D&trk=public_jobs_jserp-result_search-card", SourceID: "linkedin"},
			{Position: "Senior Frontend Developer", CompanyID: "trilogy-international-ltd", CompanyName: "Trilogy International", Location: "South San Francisco, CA", URL: "https://www.linkedin.com/jobs/view/senior-frontend-developer-at-trilogy-international-3936896077?position=3&pageNum=0&refId=fsDMYm%2BoJB2zdtWm%2FhnZ3g%3D%3D&trackingId=YzFScrnDUB9kJlZ%2FHtqdiQ%3D%3D&trk=public_jobs_jserp-result_search-card", SourceID: "linkedin"},
		}

		assert.Equal(t, want, got)
		assert.Equal(t, 0, len(errs))
		assert.Contains(t, logBufferSpy.String(), fmt.Sprintf(errMalformedCompanyLink, "fda&=+!-//"))
	})
}
