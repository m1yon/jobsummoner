package linkedin

import (
	"errors"
	"io"
	"log/slog"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/m1yon/jobsummoner/internal/models"
	"github.com/stretchr/testify/assert"
)

func TestFileReader(t *testing.T) {
	t.Run("reads the file and parses the results", func(t *testing.T) {
		fileReader := NewFileLinkedInReader("./test-helpers/li-job-listings-%v.html")

		buffer, isLastPage, err := fileReader.GetNextJobListingPage(time.Now().Add(-30 * time.Minute))
		assert.NoError(t, err)
		assert.Equal(t, false, isLastPage)

		doc, err := goquery.NewDocumentFromReader(buffer)
		assert.NoError(t, err)

		numberOfJobElements := doc.Find("body > li").Length()
		assert.Equal(t, 10, numberOfJobElements)
	})

	t.Run("handles failed opening file", func(t *testing.T) {
		fileReader := NewFileLinkedInReader("./bad-file.html")

		_, _, err := fileReader.GetNextJobListingPage(time.Now().Add(-30 * time.Minute))
		if assert.Error(t, err) {
			assert.Contains(t, err.Error(), errOpeningFile)
		}
	})
}

type stubClient struct {
	fileReader *FileLinkedInReader
}

func (m *stubClient) Get(url string) (resp *http.Response, err error) {
	if strings.Contains(url, "https://www.linkedin.com/jobs-guest/jobs/api/seeMoreJobPostings/search?f_TPR=r1800&keywords=Software+Engineer+OR+Manager&location=United+States") {
		buffer, _, err := m.fileReader.GetNextJobListingPage(time.Now().Add(-30 * time.Minute))

		if err != nil {
			return &http.Response{}, err
		}

		Body := io.NopCloser(buffer)

		return &http.Response{
			Body: Body,
		}, nil
	}

	return &http.Response{}, errors.New("bad request")
}

func NewStubClient() *stubClient {
	fileReader := NewFileLinkedInReader("./test-helpers/li-job-listings-%v.html")
	return &stubClient{fileReader}
}

func TestHttpReader(t *testing.T) {
	t.Run("hits the LinkedIn API and parses the results", func(t *testing.T) {
		stubClient := NewStubClient()
		logger := slog.New(slog.NewTextHandler(nil, nil))
		httpReader := NewHttpLinkedInReader(LinkedInReaderConfig{
			Keywords: []string{"Software Engineer", "Manager"},
			Location: "United States",
		}, stubClient, logger)

		buffer, isLastPage, err := httpReader.GetNextJobListingPage(time.Now().Add(-30 * time.Minute))
		assert.NoError(t, err)
		assert.Equal(t, false, isLastPage)

		doc, err := goquery.NewDocumentFromReader(buffer)
		assert.NoError(t, err)

		numberOfJobElements := doc.Find("body > li").Length()
		assert.Equal(t, 10, numberOfJobElements)
	})
}

func TestBuilderJobListingURL(t *testing.T) {
	tests := []struct {
		name      string
		getConfig func() LinkedInReaderConfig
		want      string
	}{
		{
			"Keywords field",
			func() LinkedInReaderConfig {
				config := LinkedInReaderConfig{Keywords: []string{"react", "typescript"}}
				return config
			},
			"?f_TPR=r1800&keywords=react+OR+typescript",
		},
		{
			"Location field",
			func() LinkedInReaderConfig {
				config := LinkedInReaderConfig{Location: "United States"}
				return config
			},
			"?f_TPR=r1800&location=United+States",
		},
		{
			"WorkTypes field",
			func() LinkedInReaderConfig {
				config := LinkedInReaderConfig{WorkTypes: []models.WorkType{models.WorkTypeRemote, models.WorkTypeOnSite}}
				return config
			},
			"?f_TPR=r1800&f_WT=2%2C1",
		},
		{
			"JobTypes field",
			func() LinkedInReaderConfig {
				config := LinkedInReaderConfig{JobTypes: []models.JobType{models.JobTypeFullTime, models.JobTypeOther}}
				return config
			},
			"?f_JT=F%2CO&f_TPR=r1800",
		},
		{
			"SalaryRange field",
			func() LinkedInReaderConfig {
				config := LinkedInReaderConfig{SalaryRange: models.SalaryRange160kPlus}
				return config
			},
			"?f_SB2=7&f_TPR=r1800",
		},
		{
			"Initial Page",
			func() LinkedInReaderConfig {
				config := LinkedInReaderConfig{
					InitialPage: 2,
				}
				return config
			},
			"?f_TPR=r1800&start=10",
		},
		{
			"All fields",
			func() LinkedInReaderConfig {
				config := LinkedInReaderConfig{
					Keywords:    []string{"go", "templ"},
					Location:    "Africa",
					WorkTypes:   []models.WorkType{models.WorkTypeHybrid},
					JobTypes:    []models.JobType{models.JobTypeFullTime, models.JobTypeOther},
					SalaryRange: models.SalaryRange200kPlus,
					InitialPage: 4,
				}
				return config
			},
			"?f_JT=F%2CO&f_SB2=9&f_TPR=r1800&f_WT=3&keywords=go+OR+templ&location=Africa&start=30",
		},
	}

	logger := slog.New(slog.NewTextHandler(nil, nil))
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := NewHttpLinkedInReader(tt.getConfig(), http.DefaultClient, logger)
			got := reader.buildJobListingURL(time.Now().Add(-30 * time.Minute))
			assert.Equal(t, linkedInBaseSearchURL+tt.want, got)
		})
	}
}
