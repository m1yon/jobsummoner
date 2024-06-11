package sqlitedb

import (
	"database/sql"
	"errors"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockSqlOpener struct {
	mock.Mock
}

func (m *MockSqlOpener) Open(driverName string, dataSourceName string) (*sql.DB, error) {
	args := m.Called(driverName, dataSourceName)
	return args.Get(0).(*sql.DB), args.Error(1)
}

func TestNewDB(t *testing.T) {
	t.Run("returns a new db connection", func(t *testing.T) {
		dataSourceName := "libsql://jobsummoner.turso.io/db"
		os.Setenv("DATABASE_URL", dataSourceName)
		defer os.Setenv("DATABASE_URL", "")

		mockOpener := new(MockSqlOpener)
		mockDb := &sql.DB{}
		mockOpener.On("Open", "libsql", dataSourceName).Return(mockDb, nil)

		db, err := NewDB(mockOpener.Open)

		mockOpener.AssertExpectations(t)
		assert.Equal(t, mockDb, db)
		assert.NoError(t, err)
	})

	t.Run("returns error when DATABASE_URL is not set", func(t *testing.T) {
		os.Setenv("DATABASE_URL", "")
		mockOpener := new(MockSqlOpener)

		_, err := NewDB(mockOpener.Open)

		mockOpener.AssertExpectations(t)
		assert.ErrorContains(t, err, ErrDatabaseURLNotSet)
	})

	t.Run("returns error when failed to open", func(t *testing.T) {
		dataSourceName := "libsql://jobsummoner.turso.io/db"
		os.Setenv("DATABASE_URL", dataSourceName)
		defer os.Setenv("DATABASE_URL", "")

		mockOpener := new(MockSqlOpener)
		mockDb := &sql.DB{}
		mockOpener.On("Open", "libsql", dataSourceName).Return(mockDb, errors.New("could not make connection"))

		_, err := NewDB(mockOpener.Open)

		mockOpener.AssertExpectations(t)
		assert.ErrorContains(t, err, ErrOpeningDB)
	})
}
