package web

import (
	"net/http"
	"testing"
	"time"

	"github.com/go-rod/rod"
)

type Driver struct {
	BaseURL string
	Client  *http.Client
	browser *rod.Browser
}

func NewWebDriver(baseURL string) (*Driver, error) {
	browser := rod.New()
	err := browser.Connect()

	if err != nil {
		return nil, err
	}

	return &Driver{
		BaseURL: "http://" + baseURL,
		Client: &http.Client{
			Timeout: 1 * time.Second,
		},
		browser: browser,
	}, nil
}

func (d Driver) SignUp(t *testing.T) (string, string, error) {
	var (
		name     = "Tom"
		email    = "tom@gmail.com"
		password = "hunter12"
	)

	page := d.browser.MustPage(d.BaseURL + "/")

	findElement(t, page, "a", "Signup").MustClick()

	nameInput := findElement(t, page, "input", "Enter your name").MustClick()
	nameInput.MustInput(name)

	emailInput := findElement(t, page, "input", "Enter your email").MustClick()
	emailInput.MustInput(email)

	passwordInput := findElement(t, page, "input", "Enter your password").MustClick()
	passwordInput.MustInput(password)

	findElement(t, page, "button", "Signup").MustClick()

	assertTextExistsInTheDocument(t, page, "Your account has successfully been created.")

	return email, password, nil
}

func (d Driver) Login(t *testing.T, email, password string) error {
	page := d.browser.MustPage(d.BaseURL + "/")

	findElement(t, page, "a", "Login").MustClick()

	emailInput := findElement(t, page, "input", "Enter your email").MustClick()
	emailInput.MustInput(email)

	passwordInput := findElement(t, page, "input", "Enter your password").MustClick()
	passwordInput.MustInput(password)

	findElement(t, page, "button", "Login").MustClick()

	assertTextExistsInTheDocument(t, page, "Login successful.")

	return nil
}

func (d Driver) AssertLoggedIn(t *testing.T) {
	page := d.browser.MustPage(d.BaseURL + "/")

	assertTextExistsInTheDocument(t, page, "Logout")
}
