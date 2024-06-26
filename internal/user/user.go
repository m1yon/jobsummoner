package user

import (
	"context"
	"database/sql"

	"github.com/pkg/errors"

	"github.com/m1yon/jobsummoner/internal/sqlitedb"
	"golang.org/x/crypto/bcrypt"
	"modernc.org/sqlite"
)

type DefaultUserService struct {
	queries *sqlitedb.Queries
}

func NewDefaultUserService(db *sql.DB) *DefaultUserService {
	queries := sqlitedb.New(db)
	return &DefaultUserService{queries}
}

func (m *DefaultUserService) Insert(name, email, password string) error {
	ctx := context.Background()

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return err
	}

	err = m.queries.CreateUser(ctx, sqlitedb.CreateUserParams{Name: name, Email: email, HashedPassword: string(hashedPassword)})

	if err != nil {
		var sqliteError *sqlite.Error

		if errors.As(err, &sqliteError) {
			if sqliteError.Code() == 2067 {
				return ErrDuplicateEmail
			}
		}

		return err
	}

	return nil
}

func (m *DefaultUserService) Authenticate(ctx context.Context, email, password string) (int, error) {
	row, err := m.queries.GetUserCredentials(ctx, email)
	if err != nil {
		return 0, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(row.HashedPassword), []byte(password))
	if err != nil {
		return 0, err
	}

	return int(row.ID), nil
}

func (m *DefaultUserService) Exists(id int) (bool, error) {
	return false, nil
}
