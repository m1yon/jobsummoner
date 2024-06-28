package models

import (
	"context"
	"strings"
	"time"

	"github.com/m1yon/jobsummoner/internal/database"
	"github.com/pkg/errors"

	"golang.org/x/crypto/bcrypt"
)

type UserModelInterface interface {
	Create(name, email, password string) error
	Authenticate(context context.Context, email, password string) (int, error)
	Exists(ctx context.Context, id int) (bool, error)
	Get(ctx context.Context, id int) (User, error)
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

func (m *UserModel) Create(name, email, password string) error {
	ctx := context.Background()

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return err
	}

	err = m.Queries.CreateUser(ctx, database.CreateUserParams{Name: name, Email: email, HashedPassword: string(hashedPassword)})

	if err != nil {
		if strings.Contains(err.Error(), DBErrUniqueConstraint) {
			return ErrDuplicateEmail
		}

		return err
	}

	return nil
}

func (m *UserModel) Authenticate(ctx context.Context, email, password string) (int, error) {
	row, err := m.Queries.GetUserCredentials(ctx, email)
	if err != nil {
		if strings.Contains(err.Error(), DBErrNoResults) {
			return 0, ErrInvalidCredentials
		}

		return 0, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(row.HashedPassword), []byte(password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return 0, ErrInvalidCredentials
		} else {
			return 0, err
		}
	}

	return int(row.ID), nil
}

func (m *UserModel) Get(ctx context.Context, id int) (User, error) {
	user, err := m.Queries.GetUser(ctx, int64(id))

	if err != nil {
		if strings.Contains(err.Error(), DBErrNoResults) {
			return User{}, ErrNoRecord
		}

		return User{}, errors.Wrap(err, "error getting user from DB")
	}

	return User{
		ID:             int(user.ID),
		Name:           user.Name,
		Email:          user.Email,
		CreatedAt:      user.CreatedAt,
		HashedPassword: []byte(user.HashedPassword),
	}, nil
}

func (m *UserModel) Exists(ctx context.Context, id int) (bool, error) {
	user, err := m.Get(ctx, id)

	if err != nil {
		if errors.Is(err, ErrNoRecord) {
			return false, nil
		}

		return false, err
	}

	return user.ID != 0, nil
}
