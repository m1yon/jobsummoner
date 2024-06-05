package components

import (
	"testing"
	"time"

	"gotest.tools/assert"
)

func TestTimeAgo(t *testing.T) {
	tests := []struct {
		name     string
		input    time.Time
		expected string
	}{
		{
			name:     "single second",
			input:    time.Now().Add(-1 * time.Second),
			expected: "1 second ago",
		},
		{
			name:     "multiple seconds",
			input:    time.Now().Add(-30 * time.Second),
			expected: "30 seconds ago",
		},
		{
			name:     "single minute",
			input:    time.Now().Add(-1 * time.Minute),
			expected: "1 minute ago",
		},
		{
			name:     "multiple minutes",
			input:    time.Now().Add(-30 * time.Minute),
			expected: "30 minutes ago",
		},
		{
			name:     "single hour",
			input:    time.Now().Add(-1 * time.Hour),
			expected: "1 hour ago",
		},
		{
			name:     "multiple hours",
			input:    time.Now().Add(-12 * time.Hour),
			expected: "12 hours ago",
		},
		{
			name:     "single day",
			input:    time.Now().Add(-24 * time.Hour),
			expected: "1 day ago",
		},
		{
			name:     "multiple days",
			input:    time.Now().Add(-336 * time.Hour),
			expected: "14 days ago",
		},
		{
			name:     "single month",
			input:    time.Now().Add(-780 * time.Hour),
			expected: "1 month ago",
		},
		{
			name:     "multiple months",
			input:    time.Now().Add(-2340 * time.Hour),
			expected: "3 months ago",
		},
		{
			name:     "single year",
			input:    time.Now().Add(-8800 * time.Hour),
			expected: "1 year ago",
		},
		{
			name:     "multiple years",
			input:    time.Now().Add(-26400 * time.Hour),
			expected: "3 years ago",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, timeAgo(tt.input), tt.expected)
		})
	}
}
