package linkedin

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/m1yon/jobsummoner"
	"github.com/pkg/errors"
)

const (
	ErrOpeningFile                 = "failed opening file"
	ErrWritingFileContentsToBuffer = "failed writing file contents to buffer"
	ErrDeterminingIfLastPage       = "failed to determine if last page"
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
	GetNextJobListingPage() (io.Reader, bool, error)
}

type HttpLinkedInReader struct {
	config LinkedInReaderConfig
	page   int
}

func (m *HttpLinkedInReader) GetNextJobListingPage() (io.Reader, bool, error) {
	url := m.buildJobListingURL()
	resp, err := http.Get(url)

	if err != nil {
		return nil, false, errors.Wrap(err, "error fetching job listing page")
	}

	if isLastPage, err := isLastJobListingPage(resp.Body); isLastPage {
		if err != nil {
			return nil, false, errors.Wrap(err, ErrDeterminingIfLastPage)
		}

		return resp.Body, true, nil
	}

	m.page++

	return resp.Body, false, nil
}

func isLastJobListingPage(reader io.Reader) (bool, error) {
	doc, err := goquery.NewDocumentFromReader(reader)

	if err != nil {
		return true, errors.Wrap(err, ErrInvalidHTML)
	}

	jobElements := doc.Find("body > li")

	if jobElements.Length() < 10 {
		return true, nil
	}

	return false, nil
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

type FileOpener interface {
	Open(name string) (*os.File, error)
	Copy(dst io.Writer, src io.Reader) (int64, error)
}

type StandardFileOpener struct{}

func (fo *StandardFileOpener) Open(name string) (*os.File, error) {
	return os.Open(name)
}

func (fo *StandardFileOpener) Copy(dst io.Writer, src io.Reader) (int64, error) {
	return io.Copy(dst, src)
}

type FileLinkedInReader struct {
	fileOpener FileOpener
	pathf      string
	page       int
}

func (m *FileLinkedInReader) GetNextJobListingPage() (io.Reader, bool, error) {
	file, err := os.Open(fmt.Sprintf(m.pathf, m.page))

	if err != nil {
		return &bytes.Buffer{}, false, errors.Wrap(err, ErrOpeningFile)
	}

	var buffer bytes.Buffer
	if _, err := m.fileOpener.Copy(&buffer, file); err != nil {
		return &bytes.Buffer{}, false, errors.Wrap(err, ErrWritingFileContentsToBuffer)
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

func NewFileLinkedInReader(pathf string, fileOpener FileOpener) *FileLinkedInReader {
	return &FileLinkedInReader{fileOpener: fileOpener, pathf: pathf, page: 1}
}
