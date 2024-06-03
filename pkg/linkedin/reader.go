package linkedin

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/m1yon/jobsummoner"
	"github.com/pkg/errors"
)

type LinkedInReaderConfig struct {
	Keywords    []string
	Location    string
	WorkTypes   []jobsummoner.WorkType
	JobTypes    []jobsummoner.JobType
	SalaryRange jobsummoner.SalaryRange
	MaxAge      time.Duration
}

const linkedInBaseSearchURL = "https://www.linkedin.com/jobs-guest/jobs/api/seeMoreJobPostings/search"

type LinkedInReader interface {
	GetNextJobListingPage() (io.Reader, error)
}

type HttpLinkedInReader struct {
	config LinkedInReaderConfig
	page   int
}

func (m *HttpLinkedInReader) GetNextJobListingPage() (io.Reader, error) {
	url := m.buildJobListingURL()
	resp, err := http.Get(url)

	if err != nil {
		return nil, errors.Wrap(err, "error fetching job listing page")
	}

	return resp.Body, nil
}

func NewHttpLinkedInReader(config LinkedInReaderConfig) *HttpLinkedInReader {
	return &HttpLinkedInReader{config, 0}
}

func (m *HttpLinkedInReader) buildJobListingURL() string {
	url, _ := url.Parse(linkedInBaseSearchURL)

	q := url.Query()
	if len(m.config.Keywords) > 0 {
		q.Set("keywords", strings.Join(m.config.Keywords, " OR "))
	}
	if m.config.Location != "" {
		q.Set("location", m.config.Location)
	}
	if len(m.config.WorkTypes) > 0 {
		q.Set("f_WT", join(m.config.WorkTypes, ","))
	}
	if len(m.config.JobTypes) > 0 {
		q.Set("f_JT", join(m.config.JobTypes, ","))
	}
	if m.config.SalaryRange != "" {
		q.Set("f_SB2", string(m.config.SalaryRange))
	}
	if m.config.MaxAge != 0.0 {
		q.Set("f_TPR", fmt.Sprintf("r%v", m.config.MaxAge.Seconds()))
	}

	url.RawQuery = q.Encode()

	return url.String()
}

func join[T ~string](input []T, sep string) string {
	slice := make([]string, len(input))
	for i, v := range input {
		slice[i] = string(v)
	}

	result := strings.Join(slice, sep)

	return result
}

type FileLinkedInReader struct {
	path string
}

func (m *FileLinkedInReader) GetNextJobListingPage() (io.Reader, error) {
	file, err := os.Open(m.path)

	if err != nil {
		return nil, errors.Wrap(err, "error opening file")
	}

	return file, nil
}

func NewFileLinkedInReader(path string) *FileLinkedInReader {
	return &FileLinkedInReader{path}
}
