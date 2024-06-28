package models

import (
	"context"
	"testing"

	"github.com/m1yon/jobsummoner/internal/database"
	_ "github.com/m1yon/jobsummoner/internal/testing"
	"github.com/stretchr/testify/assert"
)

func TestCompanies(t *testing.T) {
	companyToCreate := Company{
		ID:       "/google",
		Name:     "Google",
		Url:      "https://google.com/",
		Avatar:   "https://upload.wikimedia.org/wikipedia/commons/c/c1/Google_%22G%22_logo.svg",
		SourceID: "linkedin",
	}

	t.Run("create company and ensure it exists", func(t *testing.T) {
		ctx := context.Background()
		companies := newTestCompanyModel()

		doesCompanyExist, err := companies.Exists(ctx, companyToCreate.ID)
		assert.NoError(t, err)
		assert.Equal(t, false, doesCompanyExist, "company shouldn't exist yet")

		id, err := companies.Create(ctx, companyToCreate)
		assert.NoError(t, err)
		assert.Equal(t, companyToCreate.ID, id)

		doesCompanyExist, err = companies.Exists(ctx, id)
		assert.NoError(t, err)
		assert.Equal(t, true, doesCompanyExist, "company should exist now")
	})

	t.Run("create company and immediately get company", func(t *testing.T) {
		ctx := context.Background()
		companies := newTestCompanyModel()

		id, err := companies.Create(ctx, companyToCreate)
		assert.NoError(t, err)
		assert.Equal(t, companyToCreate.ID, id)

		company, err := companies.Get(ctx, id)
		assert.NoError(t, err)
		assert.Equal(t, companyToCreate, company, "company should exist now")
	})
}

func newTestCompanyModel() *CompanyModel {
	db, _ := database.NewInMemoryDB()

	queries := database.New(db)
	companies := &CompanyModel{queries}

	return companies
}
