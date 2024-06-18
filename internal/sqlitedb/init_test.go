package sqlitedb

import (
	"database/sql"
	"os"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

const (
	tursoFakeDataSource = "libsql://jobsummoner.turso.io/db"
)

type MockConnectionOpener struct {
	mock.Mock
}

func (m *MockConnectionOpener) Open(driverName string, dataSourceName string) (*sql.DB, error) {
	args := m.Called(driverName, dataSourceName)
	return args.Get(0).(*sql.DB), args.Error(1)
}

func TestNewFileDB(t *testing.T) {
	t.Run("returns a new file db connection", func(t *testing.T) {
		mockOpener := setupFileDBTest(t)

		db, err := NewFileDB(mockOpener)

		if assert.NoError(t, err) {
			mockOpener.AssertExpectations(t)
			assert.NotNil(t, db)
		}
	})
}

func setupFileDBTest(t *testing.T) *MockConnectionOpener {
	t.Helper()

	workingDir, err := os.Getwd()

	if err != nil {
		t.Fatal("could not get working directory")
	}

	mockOpener := new(MockConnectionOpener)
	mockDb := &sql.DB{}
	mockOpener.On("Open", "sqlite", workingDir+"/db/database.db").Return(mockDb, nil)

	return mockOpener
}

func TestNewTursoDB(t *testing.T) {
	t.Run("returns a new turso db connection", func(t *testing.T) {
		mockOpener, cleanup := setupTursoDBTest(t, tursoFakeDataSource)
		defer cleanup()

		mockOpener.On("Open", "libsql", tursoFakeDataSource).Return(&sql.DB{}, nil)

		db, err := NewTursoDB(mockOpener)

		if assert.NoError(t, err) {
			mockOpener.AssertExpectations(t)
			assert.NotNil(t, db)
		}
	})

	t.Run("returns error when DATABASE_URL is not set", func(t *testing.T) {
		mockOpener, cleanup := setupTursoDBTest(t, "")
		defer cleanup()

		_, err := NewTursoDB(mockOpener)

		mockOpener.AssertExpectations(t)
		assert.ErrorContains(t, err, ErrDatabaseURLNotSet)
	})

	t.Run("returns error when failed to open", func(t *testing.T) {
		mockOpener, cleanup := setupTursoDBTest(t, tursoFakeDataSource)
		defer cleanup()

		mockOpener.On("Open", "libsql", tursoFakeDataSource).Return(&sql.DB{}, errors.New("could not make connection"))

		_, err := NewTursoDB(mockOpener)

		mockOpener.AssertExpectations(t)
		assert.ErrorContains(t, err, ErrOpeningDB)
	})
}

func setupTursoDBTest(t *testing.T, databaseURL string) (*MockConnectionOpener, func()) {
	t.Helper()

	os.Setenv("DATABASE_URL", databaseURL)

	mockOpener := new(MockConnectionOpener)

	return mockOpener, func() {
		os.Setenv("DATABASE_URL", "")
	}
}
