package linkedin

import (
	"fmt"
	"io"
	"os"
	"testing"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/m1yon/jobsummoner"
	"github.com/stretchr/testify/assert"
)

type MockFileOpener struct{}

func (m *MockFileOpener) Open(name string) (*os.File, error) {
	return os.Open(name)
}

func (m *MockFileOpener) Copy(dst io.Writer, src io.Reader) (int64, error) {
	return 0, fmt.Errorf("simulated copy failure")
}

func TestReader(t *testing.T) {
	t.Run("reads the file", func(t *testing.T) {
		fileReader := NewFileLinkedInReader("./test-helpers/li-job-listings-%v.html", &StandardFileOpener{})

		buffer, isLastPage, err := fileReader.GetNextJobListingPage()
		assert.NoError(t, err)
		assert.Equal(t, false, isLastPage)

		doc, err := goquery.NewDocumentFromReader(buffer)
		assert.NoError(t, err)

		numberOfJobElements := doc.Find("body > li").Length()
		assert.Equal(t, 10, numberOfJobElements)
	})

	t.Run("handles failed opening file", func(t *testing.T) {
		fileReader := NewFileLinkedInReader("./bad-file.html", &StandardFileOpener{})

		_, _, err := fileReader.GetNextJobListingPage()
		if assert.Error(t, err) {
			assert.Contains(t, err.Error(), ErrOpeningFile)
		}
	})

	t.Run("handles failed writing file contents to buffer", func(t *testing.T) {
		fileReader := NewFileLinkedInReader("./test-helpers/li-job-listings-%v.html", &MockFileOpener{})

		_, _, err := fileReader.GetNextJobListingPage()
		if assert.Error(t, err) {
			assert.Contains(t, err.Error(), ErrWritingFileContentsToBuffer)
		}
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
			"All fields",
			func() LinkedInReaderConfig {
				config := LinkedInReaderConfig{
					Keywords:    []string{"go", "templ"},
					Location:    "Africa",
					WorkTypes:   []jobsummoner.WorkType{jobsummoner.WorkTypeHybrid},
					JobTypes:    []jobsummoner.JobType{jobsummoner.JobTypeFullTime, jobsummoner.JobTypeOther},
					SalaryRange: jobsummoner.SalaryRange200kPlus,
					MaxAge:      time.Hour * 12,
				}
				return config
			},
			"?f_JT=F%2CO&f_SB2=9&f_TPR=r43200&f_WT=3&keywords=go+OR+templ&location=Africa",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := NewHttpLinkedInReader(tt.getConfig())
			got := reader.buildJobListingURL()
			assert.Equal(t, linkedInBaseSearchURL+tt.want, got)
		})
	}
}
