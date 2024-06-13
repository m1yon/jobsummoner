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

func TestNewFileDB(t *testing.T) {
	t.Run("returns a new file db connection", func(t *testing.T) {
		mockOpener := setupFileDBTest(t)

		db, err := NewFileDB(mockOpener.Open)

		mockOpener.AssertExpectations(t)
		assert.NotNil(t, db)
		assert.NoError(t, err)
	})
}

func setupFileDBTest(t *testing.T) *MockSqlOpener {
	t.Helper()

	os.Setenv("LOCAL_DB", "true")
	defer os.Setenv("LOCAL_DB", "")
	workingDir, err := os.Getwd()

	if err != nil {
		t.Fatal("could not get working directory")
	}

	mockOpener := new(MockSqlOpener)
	mockDb := &sql.DB{}
	mockOpener.On("Open", "sqlite", workingDir+"/db/database.db").Return(mockDb, nil)

	return mockOpener
}

func TestNewTursoDB(t *testing.T) {
	t.Run("returns a new turso db connection", func(t *testing.T) {
		dataSourceName := "libsql://jobsummoner.turso.io/db"
		os.Setenv("DATABASE_URL", dataSourceName)
		defer os.Setenv("DATABASE_URL", "")

		mockOpener := new(MockSqlOpener)
		mockDb := &sql.DB{}
		mockOpener.On("Open", "libsql", dataSourceName).Return(mockDb, nil)

		db, err := NewTursoDB(mockOpener.Open)

		mockOpener.AssertExpectations(t)
		assert.Equal(t, mockDb, db)
		assert.NoError(t, err)
	})

	t.Run("returns error when DATABASE_URL is not set", func(t *testing.T) {
		os.Setenv("DATABASE_URL", "")
		mockOpener := new(MockSqlOpener)

		_, err := NewTursoDB(mockOpener.Open)

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

		_, err := NewTursoDB(mockOpener.Open)

		mockOpener.AssertExpectations(t)
		assert.ErrorContains(t, err, ErrOpeningDB)
	})
}
