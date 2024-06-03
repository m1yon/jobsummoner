package linkedin

import (
	"testing"
	"time"

	"github.com/m1yon/jobsummoner"
	"github.com/stretchr/testify/assert"
)

func TestLinkedInURLBuilder(t *testing.T) {
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
			got := reader.buildLinkedInJobsURL()
			assert.Equal(t, linkedInBaseSearchURL+tt.want, got)
		})
	}
}
