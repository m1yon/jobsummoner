package sqlitedb

import (
	"context"
	"database/sql"

	"github.com/m1yon/jobsummoner"
	"github.com/pkg/errors"
)

type SqliteCompanyRepository struct {
	queries *Queries
}

func NewSqliteCompanyRepository(db *sql.DB) *SqliteCompanyRepository {
	queries := New(db)
	return &SqliteCompanyRepository{queries}
}

func (s *SqliteCompanyRepository) DoesCompanyExist(ctx context.Context, id string) (bool, error) {
	company, err := s.queries.GetCompany(ctx, id)

	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return false, nil
		}

		return false, errors.Wrap(err, "error getting company from DB")
	}

	return company.ID != "", nil
}

func (s *SqliteCompanyRepository) CreateCompany(ctx context.Context, company jobsummoner.Company) (string, error) {
	err := s.queries.CreateCompany(ctx, CreateCompanyParams{
		ID:       company.ID,
		Name:     company.Name,
		Url:      company.Url,
		Avatar:   sql.NullString{String: company.Avatar, Valid: company.Avatar != ""},
		SourceID: company.SourceID,
	})

	if err != nil {
		return "", errors.Wrap(err, "error creating company in DB")
	}

	return company.ID, nil
}

func NewInMemorySqliteCompanyRepository(db *sql.DB) *SqliteCompanyRepository {
	queries := New(db)
	return &SqliteCompanyRepository{queries}
}
