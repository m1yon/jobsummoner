package linkedin

import (
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/m1yon/jobsummoner"
	"github.com/stretchr/testify/assert"
)

func TestFileReader(t *testing.T) {
	t.Run("reads the file and parses the results", func(t *testing.T) {
		fileReader := NewFileLinkedInReader("./test-helpers/li-job-listings-%v.html")

		buffer, isLastPage, err := fileReader.GetNextJobListingPage()
		assert.NoError(t, err)
		assert.Equal(t, false, isLastPage)

		doc, err := goquery.NewDocumentFromReader(buffer)
		assert.NoError(t, err)

		numberOfJobElements := doc.Find("body > li").Length()
		assert.Equal(t, 10, numberOfJobElements)
	})

	t.Run("handles failed opening file", func(t *testing.T) {
		fileReader := NewFileLinkedInReader("./bad-file.html")

		_, _, err := fileReader.GetNextJobListingPage()
		if assert.Error(t, err) {
			assert.Contains(t, err.Error(), errOpeningFile)
		}
	})
}

type stubClient struct {
	fileReader *FileLinkedInReader
}

func (m *stubClient) Get(url string) (resp *http.Response, err error) {
	if strings.Contains(url, "https://www.linkedin.com/jobs-guest/jobs/api/seeMoreJobPostings/search?f_TPR=r14400&keywords=Software+Engineer+OR+Manager&location=United+States") {
		buffer, _, err := m.fileReader.GetNextJobListingPage()

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
		httpReader := NewHttpLinkedInReader(LinkedInReaderConfig{
			Keywords: []string{"Software Engineer", "Manager"},
			Location: "United States",
			MaxAge:   time.Hour * 4,
		}, stubClient)

		buffer, isLastPage, err := httpReader.GetNextJobListingPage()
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
			"?keywords=react+OR+typescript",
		},
		{
			"Location field",
			func() LinkedInReaderConfig {
				config := LinkedInReaderConfig{Location: "United States"}
				return config
			},
			"?location=United+States",
		},
		{
			"WorkTypes field",
			func() LinkedInReaderConfig {
				config := LinkedInReaderConfig{WorkTypes: []jobsummoner.WorkType{jobsummoner.WorkTypeRemote, jobsummoner.WorkTypeOnSite}}
				return config
			},
			"?f_WT=2%2C1",
		},
		{
			"JobTypes field",
			func() LinkedInReaderConfig {
				config := LinkedInReaderConfig{JobTypes: []jobsummoner.JobType{jobsummoner.JobTypeFullTime, jobsummoner.JobTypeOther}}
				return config
			},
			"?f_JT=F%2CO",
		},
		{
			"SalaryRange field",
			func() LinkedInReaderConfig {
				config := LinkedInReaderConfig{SalaryRange: jobsummoner.SalaryRange160kPlus}
				return config
			},
			"?f_SB2=7",
		},
		{
			"MaxAge field",
			func() LinkedInReaderConfig {
				config := LinkedInReaderConfig{MaxAge: time.Hour * 24}
				return config
			},
			"?f_TPR=r86400",
		},
		{
			"Initial Page",
			func() LinkedInReaderConfig {
				config := LinkedInReaderConfig{
					InitialPage: 2,
				}
				return config
			},
			"?start=10",
		},
		{
			"All fields",
			func() LinkedInReaderConfig {
				config := LinkedInReaderConfig{
					Keywords:    []string{"go", "templ"},
					Location:    "Africa",
					WorkTypes:   []jobsummoner.WorkType{jobsummoner.WorkTypeHybrid},
					JobTypes:    []jobsummoner.JobType{jobsummoner.JobTypeFullTime, jobsummoner.JobTypeOther},
					SalaryRange: jobsummoner.SalaryRange200kPlus,
					MaxAge:      time.Hour * 12,
					InitialPage: 4,
				}
				return config
			},
			"?f_JT=F%2CO&f_SB2=9&f_TPR=r43200&f_WT=3&keywords=go+OR+templ&location=Africa&start=30",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := NewHttpLinkedInReader(tt.getConfig(), http.DefaultClient)
			got := reader.buildJobListingURL()
			assert.Equal(t, linkedInBaseSearchURL+tt.want, got)
		})
	}
}
