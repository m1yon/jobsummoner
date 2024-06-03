package linkedin

import (
	"io"
	"os"

	"github.com/pkg/errors"
)

type MockLinkedInReader struct {
	path string
}

func (m *MockLinkedInReader) GetJobListingPage(itemOffset int) (io.Reader, error) {
	file, err := os.Open(m.path)

	if err != nil {
		return nil, errors.Wrap(err, "error opening file")
	}

	return file, nil
}

func NewMockLinkedInReader(path string) *MockLinkedInReader {
	return &MockLinkedInReader{path}
}
