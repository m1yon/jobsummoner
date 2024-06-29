package specifications

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type AuthDriver interface {
	SignUp(t *testing.T) (email, password string, err error)
	Login(t *testing.T, email, password string) error
	AssertLoggedIn(t *testing.T)
}

func AuthSpecification(t *testing.T, app AuthDriver) {
	t.Run("can signup and login", func(t *testing.T) {
		email, password, err := app.SignUp(t)
		assert.NoError(t, err)

		err = app.Login(t, email, password)
		assert.NoError(t, err)

		app.AssertLoggedIn(t)
	})
}
