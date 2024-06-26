package main

import (
	"context"
	"errors"
	"net/http"

	"github.com/m1yon/jobsummoner/internal/components"
	"github.com/m1yon/jobsummoner/internal/models"
	"github.com/m1yon/jobsummoner/pkg/validator"
)

func (s *Server) getHomepageHandler(w http.ResponseWriter, r *http.Request) {
	jobs, err := s.jobService.GetJobs(r.Context())

	if err != nil {
		s.serverError(w, r, err)
	}

	m := s.NewHomepageViewModel(r, jobs)
	component := components.Homepage(m)
	err = s.Render(component, context.Background(), w)

	if err != nil {
		s.serverError(w, r, err)
	}
}

func (s *Server) userSignup(w http.ResponseWriter, r *http.Request) {
	component := components.SignupPage()
	err := s.Render(component, context.Background(), w)

	if err != nil {
		s.serverError(w, r, err)
	}
}

func (s *Server) userSignupPost(w http.ResponseWriter, r *http.Request) {
	var form components.UserSignupForm

	err := s.decodePostForm(r, &form)
	if err != nil {
		s.clientError(w, http.StatusBadRequest)
		return
	}

	form.CheckField(validator.NotBlank(form.Name), "name", "This field cannot be blank")
	form.CheckField(validator.NotBlank(form.Email), "email", "This field cannot be blank")
	form.CheckField(validator.Matches(form.Email, validator.EmailRX), "email", "This field must be a valid email address")
	form.CheckField(validator.NotBlank(form.Password), "password", "This field cannot be blank")
	form.CheckField(validator.MinChars(form.Password, 8), "password", "This field must be at least 8 characters long")

	if !form.Valid() {
		s.render(w, r, http.StatusOK, components.SignupForm(form))
		return
	}

	err = s.users.Insert(form.Name, form.Email, form.Password)
	if err != nil {
		if errors.Is(err, models.ErrDuplicateEmail) {
			form.AddFieldError("email", "Email address is already in use")
			s.render(w, r, http.StatusOK, components.SignupForm(form))
			return
		}

		s.serverError(w, r, err)
		return
	}

	s.sessionManager.Put(r.Context(), "flash", "Your account has successfully been created.")

	w.Header().Set("HX-Redirect", "/")
	w.WriteHeader(http.StatusSeeOther)
}

func (s *Server) userLogin(w http.ResponseWriter, r *http.Request) {
	component := components.LoginPage()
	err := s.Render(component, context.Background(), w)

	if err != nil {
		s.serverError(w, r, err)
	}
}

func (s *Server) userLoginPost(w http.ResponseWriter, r *http.Request) {
	var form components.UserLoginForm

	err := s.decodePostForm(r, &form)
	if err != nil {
		s.clientError(w, http.StatusBadRequest)
		return
	}

	form.CheckField(validator.NotBlank(form.Email), "email", "This field cannot be blank")
	form.CheckField(validator.NotBlank(form.Password), "password", "This field cannot be blank")

	if !form.Valid() {
		s.render(w, r, http.StatusOK, components.LoginForm(form))
		return
	}

	id, err := s.users.Authenticate(r.Context(), form.Email, form.Password)
	if err != nil {
		form.AddNonFieldError("Email or password is incorrect")
		s.render(w, r, http.StatusOK, components.LoginForm(form))
		return
	}

	s.sessionManager.Put(r.Context(), "authenticatedUserID", id)
	s.sessionManager.Put(r.Context(), "flash", "Login successful.")

	w.Header().Set("HX-Redirect", "/")
	w.WriteHeader(http.StatusSeeOther)
}

func (s *Server) userLogoutPost(w http.ResponseWriter, r *http.Request) {
	err := s.sessionManager.RenewToken(r.Context())
	if err != nil {
		s.serverError(w, r, err)
		return
	}

	s.sessionManager.Remove(r.Context(), "authenticatedUserID")
	s.sessionManager.Put(r.Context(), "flash", "You've been logged out successfully!")

	w.Header().Set("HX-Redirect", "/")
	w.WriteHeader(http.StatusSeeOther)
}
