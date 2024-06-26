package linkedin

import (
	"bytes"
	"fmt"
	"io"
	"log/slog"
	"math"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/m1yon/jobsummoner/internal/models"
	"github.com/pkg/errors"
)

const (
	errOpeningFile                 = "failed opening file"
	errWritingFileContentsToBuffer = "failed writing file contents to buffer"
	errDeterminingIfLastPage       = "failed to determine if last page"
)

type LinkedInReaderConfig struct {
	Keywords    []string
	Location    string
	WorkTypes   []models.WorkType
	JobTypes    []models.JobType
	SalaryRange models.SalaryRange
	InitialPage int
}

const linkedInBaseSearchURL = "https://www.linkedin.com/jobs-guest/jobs/api/seeMoreJobPostings/search"

type LinkedInReader interface {
	GetNextJobListingPage(lastScraped time.Time) (io.Reader, bool, error)
}

type HttpLinkedInReader struct {
	config LinkedInReaderConfig
	page   int
	client httpGetter
	logger *slog.Logger
}

func NewHttpLinkedInReader(config LinkedInReaderConfig, client httpGetter, logger *slog.Logger) *HttpLinkedInReader {
	return &HttpLinkedInReader{config, config.InitialPage, client, logger}
}

func (m *HttpLinkedInReader) GetNextJobListingPage(lastScraped time.Time) (io.Reader, bool, error) {
	url := m.buildJobListingURL(lastScraped)
	resp, err := m.client.Get(url)

	m.logger.Debug("reading linkedin url", slog.String("url", url))

	if err != nil {
		return &bytes.Buffer{}, false, errors.Wrap(err, "failed fetching job listings page")
	}

	var buffer bytes.Buffer
	if _, err := io.Copy(&buffer, resp.Body); err != nil {
		return &bytes.Buffer{}, false, errors.Wrap(err, "failed writing response body contents to buffer")
	}

	if isLastPage, err := isLastJobListingPage(bytes.NewBuffer(buffer.Bytes())); isLastPage {
		if err != nil {
			return &bytes.Buffer{}, false, errors.Wrap(err, "error determining if job listing page is last")
		}

		return &buffer, true, nil
	}

	m.page++

	return &buffer, false, nil
}

func isLastJobListingPage(reader io.Reader) (bool, error) {
	doc, err := goquery.NewDocumentFromReader(reader)

	if err != nil {
		return true, errors.Wrap(err, errInvalidHTML)
	}

	jobElements := doc.Find("body > li")

	if jobElements.Length() < 10 {
		return true, nil
	}

	return false, nil
}

type httpGetter interface {
	Get(url string) (resp *http.Response, err error)
}

func (m *HttpLinkedInReader) buildJobListingURL(timeSinceLastScrape time.Time) string {
	url, _ := url.Parse(linkedInBaseSearchURL)

	q := url.Query()
	q.Set("f_TPR", fmt.Sprintf("r%v", math.Floor(time.Since(timeSinceLastScrape).Seconds())))
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
	if (m.page - 1) > 0 {
		q.Set("start", fmt.Sprintf("%v", (m.page-1)*10))
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
	pathf string
	page  int
}

func NewFileLinkedInReader(pathf string) *FileLinkedInReader {
	return &FileLinkedInReader{pathf: pathf, page: 1}
}

func (m *FileLinkedInReader) GetNextJobListingPage(lastScraped time.Time) (io.Reader, bool, error) {
	file, err := os.Open(fmt.Sprintf(m.pathf, m.page))

	if err != nil {
		return &bytes.Buffer{}, false, errors.Wrap(err, errOpeningFile)
	}

	var buffer bytes.Buffer
	if _, err := io.Copy(&buffer, file); err != nil {
		return &bytes.Buffer{}, false, errors.Wrap(err, errWritingFileContentsToBuffer)
	}

	if isLastPage, err := isLastJobListingPage(bytes.NewBuffer(buffer.Bytes())); isLastPage {
		if err != nil {
			return &bytes.Buffer{}, false, errors.Wrap(err, "error determining if job listing page is last")
		}

		return &buffer, true, nil
	}

	m.page++

	return &buffer, false, nil
}
