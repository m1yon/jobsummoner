package sqlitedb

import (
	"context"
	"database/sql"
	"testing"

	"github.com/m1yon/jobsummoner"
	"github.com/stretchr/testify/assert"
)

func NewInMemorySqliteCompanyRepository(db *sql.DB, jobRepository jobsummoner.JobRepository) *SqliteCompanyRepository {
	queries := New(db)
	return &SqliteCompanyRepository{jobRepository, queries}
}

func TestCompanyRepository(t *testing.T) {
	t.Run("add company and ensure it exists", func(t *testing.T) {
		ctx := context.Background()
		db := NewTestDB()
		jobRepository := NewInMemorySqliteJobRepository(db)
		companyRepository := NewInMemorySqliteCompanyRepository(db, jobRepository)

		companyToAdd := jobsummoner.Company{
			ID:       "/google",
			Name:     "Google",
			Url:      "https://google.com/",
			Avatar:   "https://upload.wikimedia.org/wikipedia/commons/c/c1/Google_%22G%22_logo.svg",
			SourceID: "1",
		}

		doesCompanyExist, err := companyRepository.DoesCompanyExist(ctx, companyToAdd.ID)
		assert.NoError(t, err)
		assert.Equal(t, false, doesCompanyExist)

		id, err := companyRepository.AddCompany(ctx, companyToAdd)
		assert.NoError(t, err)
		assert.Equal(t, companyToAdd.ID, id)

		doesCompanyExist, err = companyRepository.DoesCompanyExist(ctx, id)
		assert.NoError(t, err)
		assert.Equal(t, true, doesCompanyExist)
	})
}
