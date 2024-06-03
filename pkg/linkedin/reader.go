package linkedin

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/m1yon/jobsummoner"
	"github.com/pkg/errors"
)

const linkedInBaseSearchURL = "https://www.linkedin.com/jobs-guest/jobs/api/seeMoreJobPostings/search"

type LinkedInReader interface {
	GetJobListingPage(itemOffset int) (io.Reader, error)
}

type HttpLinkedInReader struct {
}

func (m *HttpLinkedInReader) GetJobListingPage(itemOffset int) (io.Reader, error) {
	url := BuildLinkedInJobsURL(BuildLinkedInJobsURLArgs{
		Keywords: []string{"go"},
		Location: "United States",
		MaxAge:   time.Hour * 12,
	})
	resp, err := http.Get(url)

	if err != nil {
		return nil, errors.Wrap(err, "error fetching job listing page")
	}

	return resp.Body, nil
}

func NewHttpLinkedInReader() *HttpLinkedInReader {
	return &HttpLinkedInReader{}
}

type BuildLinkedInJobsURLArgs struct {
	Keywords    []string
	Location    string
	WorkTypes   []jobsummoner.WorkType
	JobTypes    []jobsummoner.JobType
	SalaryRange jobsummoner.SalaryRange
	MaxAge      time.Duration
}

func BuildLinkedInJobsURL(args BuildLinkedInJobsURLArgs) string {
	url, _ := url.Parse(linkedInBaseSearchURL)

	q := url.Query()
	if len(args.Keywords) > 0 {
		q.Set("keywords", strings.Join(args.Keywords, " OR "))
	}
	if args.Location != "" {
		q.Set("location", args.Location)
	}
	if len(args.WorkTypes) > 0 {
		q.Set("f_WT", join(args.WorkTypes, ","))
	}
	if len(args.JobTypes) > 0 {
		q.Set("f_JT", join(args.JobTypes, ","))
	}
	if args.SalaryRange != "" {
		q.Set("f_SB2", string(args.SalaryRange))
	}
	if args.MaxAge != 0.0 {
		q.Set("f_TPR", fmt.Sprintf("r%v", args.MaxAge.Seconds()))
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
