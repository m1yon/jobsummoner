package company

import (
	"context"
	"testing"

	"github.com/m1yon/jobsummoner"
	"github.com/m1yon/jobsummoner/internal/sqlitedb"
	"github.com/stretchr/testify/assert"
)

func TestSqliteCompanyService(t *testing.T) {
	t.Run("create company and immediately get created company", func(t *testing.T) {
		ctx := context.Background()
		db := sqlitedb.NewTestDB()
		companyRepository := sqlitedb.NewInMemorySqliteCompanyRepository(db)
		companyService := NewDefaultCompanyService(companyRepository)

		companyToCreate := jobsummoner.Company{
			ID:       "/google",
			Name:     "Google",
			Url:      "https://google.com/",
			Avatar:   "https://upload.wikimedia.org/wikipedia/commons/c/c1/Google_%22G%22_logo.svg",
			SourceID: "1",
		}

		doesCompanyExist, err := companyService.DoesCompanyExist(ctx, companyToCreate.ID)
		assert.NoError(t, err)
		assert.Equal(t, false, doesCompanyExist, "company shouldn't exist yet")

		id, err := companyRepository.AddCompany(ctx, companyToCreate)
		assert.NoError(t, err)
		assert.Equal(t, companyToCreate.ID, id)

		doesCompanyExist, err = companyRepository.DoesCompanyExist(ctx, id)
		assert.NoError(t, err)
		assert.Equal(t, true, doesCompanyExist, "company should exist now")
	})
}
