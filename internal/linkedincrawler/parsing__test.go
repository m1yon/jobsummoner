package linkedincrawler

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestParseRelativeTime(t *testing.T) {
	now := time.Now()
	tests := map[string]struct {
		inputTime   time.Time
		inputString string
		wantTime    time.Time
		wantError   error
	}{
		"x minute ago":  {now, "1 minute ago", now.Add(-(time.Minute * 1)), nil},
		"x minutes ago": {now, "5 minutes ago", now.Add(-(time.Minute * 5)), nil},
		"x hour ago":    {now, "1 hour ago", now.Add(-(time.Hour * 1)), nil},
		"x hours ago":   {now, "3 hours ago", now.Add(-(time.Hour * 3)), nil},
		"x day ago":     {now, "1 day ago", now.Add(-(time.Hour * 24)), nil},
		"x days ago":    {now, "2 days ago", now.Add(-(time.Hour * 48)), nil},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			time, err := parseRelativeTime(tt.inputTime, tt.inputString)

			assert.Equal(t, time, tt.wantTime, "they should be equal")
			assert.Equal(t, err, tt.wantError, "they should be equal")
		})
	}
}
