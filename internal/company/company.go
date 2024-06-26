package company

import (
	"context"
	"database/sql"

	"github.com/m1yon/jobsummoner"
	"github.com/m1yon/jobsummoner/internal/sqlitedb"
	"github.com/pkg/errors"
)

type DefaultCompanyService struct {
	Queries *sqlitedb.Queries
}

func NewDefaultCompanyService(queries *sqlitedb.Queries) *DefaultCompanyService {
	return &DefaultCompanyService{queries}
}

func (c *DefaultCompanyService) CreateCompany(ctx context.Context, company jobsummoner.Company) (string, error) {
	err := c.Queries.CreateCompany(ctx, sqlitedb.CreateCompanyParams{
		ID:       company.ID,
		Name:     company.Name,
		Url:      company.Url,
		Avatar:   sql.NullString{String: company.Avatar, Valid: company.Avatar != ""},
		SourceID: company.SourceID,
	})

	if err != nil {
		return "", errors.Wrap(err, "error creating company in company service")
	}

	return company.ID, nil
}

func (c *DefaultCompanyService) DoesCompanyExist(ctx context.Context, id string) (bool, error) {
	company, err := c.GetCompany(ctx, id)

	if err != nil {
		return false, errors.Wrap(err, "error finding company in company service")
	}

	return company.ID != "", nil
}

func (c *DefaultCompanyService) GetCompany(ctx context.Context, id string) (jobsummoner.Company, error) {
	company, err := c.Queries.GetCompany(ctx, id)

	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return jobsummoner.Company{}, nil
		}

		return jobsummoner.Company{}, errors.Wrap(err, "error getting company from DB")
	}

	return jobsummoner.Company{
		ID:       company.ID,
		Name:     company.Name,
		Url:      company.Url,
		Avatar:   company.Avatar.String,
		SourceID: company.SourceID,
	}, nil
}
