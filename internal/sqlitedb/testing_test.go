package sqlitedb

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateTestDB(t *testing.T) {
	_ = NewTestDB()
	assert.Equal(t, 1, 1)
}
