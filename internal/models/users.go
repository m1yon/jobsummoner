package models

import (
	"context"
	"time"

	"github.com/m1yon/jobsummoner/internal/database"
	"github.com/pkg/errors"

	"golang.org/x/crypto/bcrypt"
	"modernc.org/sqlite"
)

type UserModelInterface interface {
	Insert(name, email, password string) error
	Authenticate(context context.Context, email, password string) (int, error)
	Exists(id int) (bool, error)
}

type User struct {
	ID             int
	Name           string
	Email          string
	HashedPassword []byte
	CreatedAt      time.Time
}
type UserModel struct {
	Queries *database.Queries
}

func (m *UserModel) Insert(name, email, password string) error {
	ctx := context.Background()

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return err
	}

	err = m.Queries.CreateUser(ctx, database.CreateUserParams{Name: name, Email: email, HashedPassword: string(hashedPassword)})

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

func (m *UserModel) Authenticate(ctx context.Context, email, password string) (int, error) {
	row, err := m.Queries.GetUserCredentials(ctx, email)
	if err != nil {
		return 0, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(row.HashedPassword), []byte(password))
	if err != nil {
		return 0, err
	}

	return int(row.ID), nil
}

func (m *UserModel) Exists(id int) (bool, error) {
	return false, nil
}
