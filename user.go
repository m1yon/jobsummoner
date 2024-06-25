package jobsummoner

import (
	"context"
	"time"
)

type User struct {
	ID             int
	Name           string
	Email          string
	HashedPassword []byte
	CreatedAt      time.Time
}

type UserService interface {
	Insert(name, email, password string) error
	Authenticate(context context.Context, email, password string) (int, error)
	Exists(id int) (bool, error)
}
