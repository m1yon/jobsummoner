package models

import (
	"context"
	"database/sql"

	"github.com/m1yon/jobsummoner/internal/database"
	"github.com/pkg/errors"
)

type CompanyModelInterface interface {
	CreateCompany(ctx context.Context, company Company) (string, error)
	DoesCompanyExist(ctx context.Context, id string) (bool, error)
	GetCompany(ctx context.Context, id string) (Company, error)
}

type Company struct {
	ID       string
	Name     string
	Url      string
	Avatar   string
	SourceID string
}

type CompanyModel struct {
	Queries *database.Queries
}

func (m *CompanyModel) CreateCompany(ctx context.Context, company Company) (string, error) {
	err := m.Queries.CreateCompany(ctx, database.CreateCompanyParams{
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

func (m *CompanyModel) DoesCompanyExist(ctx context.Context, id string) (bool, error) {
	company, err := m.GetCompany(ctx, id)

	if err != nil {
		return false, errors.Wrap(err, "error finding company in company service")
	}

	return company.ID != "", nil
}

func (m *CompanyModel) GetCompany(ctx context.Context, id string) (Company, error) {
	company, err := m.Queries.GetCompany(ctx, id)

	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return Company{}, nil
		}

		return Company{}, errors.Wrap(err, "error getting company from DB")
	}

	return Company{
		ID:       company.ID,
		Name:     company.Name,
		Url:      company.Url,
		Avatar:   company.Avatar.String,
		SourceID: company.SourceID,
	}, nil
}
