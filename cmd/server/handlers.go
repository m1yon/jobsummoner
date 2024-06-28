package main

import (
	"context"
	"errors"
	"net/http"

	"github.com/m1yon/jobsummoner/internal/components"
	"github.com/m1yon/jobsummoner/internal/models"
	"github.com/m1yon/jobsummoner/pkg/validator"
)

func (app *application) getHomepageHandler(w http.ResponseWriter, r *http.Request) {
	jobs, err := app.jobs.GetMany(r.Context())

	if err != nil {
		app.serverError(w, r, err)
	}

	m := app.NewHomepageViewModel(r, jobs)
	component := components.Homepage(m)
	err = app.Render(component, context.Background(), w)

	if err != nil {
		app.serverError(w, r, err)
	}
}

func (app *application) userSignup(w http.ResponseWriter, r *http.Request) {
	component := components.SignupPage()
	err := app.Render(component, context.Background(), w)

	if err != nil {
		app.serverError(w, r, err)
	}
}

func (app *application) userSignupPost(w http.ResponseWriter, r *http.Request) {
	var form components.UserSignupForm

	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form.CheckField(validator.NotBlank(form.Name), "name", "This field cannot be blank")
	form.CheckField(validator.NotBlank(form.Email), "email", "This field cannot be blank")
	form.CheckField(validator.Matches(form.Email, validator.EmailRX), "email", "This field must be a valid email address")
	form.CheckField(validator.NotBlank(form.Password), "password", "This field cannot be blank")
	form.CheckField(validator.MinChars(form.Password, 8), "password", "This field must be at least 8 characters long")

	if !form.Valid() {
		app.render(w, r, http.StatusOK, components.SignupForm(form))
		return
	}

	err = app.users.Create(form.Name, form.Email, form.Password)
	if err != nil {
		if errors.Is(err, models.ErrDuplicateEmail) {
			form.AddFieldError("email", "Email address is already in use")
			app.render(w, r, http.StatusOK, components.SignupForm(form))
			return
		}

		app.serverError(w, r, err)
		return
	}

	app.sessionManager.Put(r.Context(), "flash", "Your account has successfully been created.")

	w.Header().Set("HX-Redirect", "/")
	w.WriteHeader(http.StatusSeeOther)
}

func (app *application) userLogin(w http.ResponseWriter, r *http.Request) {
	component := components.LoginPage()
	err := app.Render(component, context.Background(), w)

	if err != nil {
		app.serverError(w, r, err)
	}
}

func (app *application) userLoginPost(w http.ResponseWriter, r *http.Request) {
	var form components.UserLoginForm

	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form.CheckField(validator.NotBlank(form.Email), "email", "This field cannot be blank")
	form.CheckField(validator.NotBlank(form.Password), "password", "This field cannot be blank")

	if !form.Valid() {
		app.render(w, r, http.StatusOK, components.LoginForm(form))
		return
	}

	id, err := app.users.Authenticate(r.Context(), form.Email, form.Password)
	if err != nil {
		form.AddNonFieldError("Email or password is incorrect")
		app.render(w, r, http.StatusOK, components.LoginForm(form))
		return
	}

	app.sessionManager.Put(r.Context(), "authenticatedUserID", id)
	app.sessionManager.Put(r.Context(), "flash", "Login successful.")

	w.Header().Set("HX-Redirect", "/")
	w.WriteHeader(http.StatusSeeOther)
}

func (app *application) userLogoutPost(w http.ResponseWriter, r *http.Request) {
	err := app.sessionManager.RenewToken(r.Context())
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	app.sessionManager.Remove(r.Context(), "authenticatedUserID")
	app.sessionManager.Put(r.Context(), "flash", "You've been logged out successfully!")

	w.Header().Set("HX-Redirect", "/")
	w.WriteHeader(http.StatusSeeOther)
}
