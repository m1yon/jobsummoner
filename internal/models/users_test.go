package models

import (
	"context"
	"testing"

	"github.com/m1yon/jobsummoner/internal/database"
	"github.com/stretchr/testify/assert"
)

func TestUser(t *testing.T) {
	userToCreate := User{
		Name:  "Michael",
		Email: "fake@gmail.com",
	}

	t.Run("create user and immediately ensure user exists", func(t *testing.T) {
		users := newTestUserModel(t)

		err := users.Create(userToCreate.Name, userToCreate.Email, "hunter12")
		assert.NoError(t, err)

		assertUserExist(t, users, 1)
	})

	t.Run("create user and immediately get created user", func(t *testing.T) {
		users := newTestUserModel(t)

		err := users.Create(userToCreate.Name, userToCreate.Email, "hunter12")
		assert.NoError(t, err)

		user, err := users.Get(context.Background(), 1)
		assert.NoError(t, err)

		assertUsersAreEqual(t, userToCreate, user)
	})

	t.Run("get non-existing user", func(t *testing.T) {
		users := newTestUserModel(t)

		_, err := users.Get(context.Background(), 2)
		assert.ErrorIs(t, err, ErrNoRecord)
	})

	t.Run("prevent dupliate emails", func(t *testing.T) {
		users := newTestUserModel(t)

		err := users.Create(userToCreate.Name, userToCreate.Email, "hunter12")
		assert.NoError(t, err)
		err = users.Create(userToCreate.Name, userToCreate.Email, "hunter12")
		assert.ErrorIs(t, err, ErrDuplicateEmail)
	})

	t.Run("successfully authenticate", func(t *testing.T) {
		users := newTestUserModel(t)

		err := users.Create(userToCreate.Name, userToCreate.Email, "hunter12")
		assert.NoError(t, err)

		id, err := users.Authenticate(context.Background(), userToCreate.Email, "hunter12")
		assert.NoError(t, err)
		assert.Equal(t, 1, id)
	})

	t.Run("invalid username", func(t *testing.T) {
		users := newTestUserModel(t)

		err := users.Create(userToCreate.Name, "badusername", "hunter12")
		assert.NoError(t, err)

		id, err := users.Authenticate(context.Background(), userToCreate.Email, "hunter12")
		assert.ErrorIs(t, err, ErrInvalidCredentials)
		assert.Equal(t, 0, id)
	})

	t.Run("invalid password", func(t *testing.T) {
		users := newTestUserModel(t)

		err := users.Create(userToCreate.Name, userToCreate.Email, "hunter12")
		assert.NoError(t, err)

		id, err := users.Authenticate(context.Background(), userToCreate.Email, "wrongpassword")
		assert.ErrorIs(t, err, ErrInvalidCredentials)
		assert.Equal(t, 0, id)
	})

	t.Run("check if non-existing user exists", func(t *testing.T) {
		users := newTestUserModel(t)

		exists, err := users.Exists(context.Background(), 2)
		assert.NoError(t, err)
		assert.Equal(t, false, exists)
	})
}

func newTestUserModel(t *testing.T) *UserModel {
	db, err := database.NewInMemoryDB()

	if err != nil {
		t.Fatal(err)
	}

	queries := database.New(db)
	users := &UserModel{queries}

	return users
}

func assertUsersAreEqual(t *testing.T, expectedUser, actualUser User) {
	assert.Equal(t, expectedUser.Name, actualUser.Name)
	assert.Equal(t, expectedUser.Email, actualUser.Email)
}

func assertUserExist(t *testing.T, users UserModelInterface, userID int) {
	t.Helper()

	doesUserExist, err := users.Exists(context.Background(), userID)
	assert.NoError(t, err)
	assert.Equal(t, true, doesUserExist)
}
