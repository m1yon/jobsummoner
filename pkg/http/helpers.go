package http

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/a-h/templ"
	"github.com/go-playground/form/v4"
)

func (s *Server) serverError(w http.ResponseWriter, r *http.Request, err error) {
	var (
		method = r.Method
		uri    = r.URL.RequestURI()
	)

	s.logger.Error(err.Error(), "method", method, "uri", uri)
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func (s *Server) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

func (s *Server) decodePostForm(r *http.Request, dst any) error {
	err := r.ParseForm()
	if err != nil {
		return err
	}

	err = s.formDecoder.Decode(dst, r.PostForm)
	if err != nil {
		var invalidDecoderError *form.InvalidDecoderError

		if errors.As(err, &invalidDecoderError) {
			panic(err)
		}
		return err
	}

	return nil
}

func (s *Server) render(w http.ResponseWriter, r *http.Request, status int, component templ.Component) {
	w.WriteHeader(status)

	err := s.Render(component, context.Background(), w)
	if err != nil {
		s.serverError(w, r, err)
		return
	}
}

func timeAgo(from time.Time) string {
	now := time.Now()
	diff := now.Sub(from)

	if diff < time.Minute {
		if int(diff.Seconds()) == 1 {
			return fmt.Sprintf("%d second ago", int(diff.Seconds()))
		}
		return fmt.Sprintf("%d seconds ago", int(diff.Seconds()))
	} else if diff < time.Hour {
		if int(diff.Minutes()) == 1 {
			return fmt.Sprintf("%d minute ago", int(diff.Minutes()))
		}
		return fmt.Sprintf("%d minutes ago", int(diff.Minutes()))
	} else if diff < time.Hour*24 {
		if int(diff.Hours()) == 1 {
			return fmt.Sprintf("%d hour ago", int(diff.Hours()))
		}
		return fmt.Sprintf("%d hours ago", int(diff.Hours()))
	} else if diff < time.Hour*24*30 {
		days := diff / (time.Hour * 24)
		if days == 1 {
			return fmt.Sprintf("%d day ago", int(days))
		}
		return fmt.Sprintf("%d days ago", int(days))
	} else if diff < time.Hour*24*365 {
		months := diff / (time.Hour * 24 * 30)
		if months == 1 {
			return fmt.Sprintf("%d month ago", int(months))
		}
		return fmt.Sprintf("%d months ago", int(months))
	}
	years := diff / (time.Hour * 24 * 365)
	if years == 1 {
		return fmt.Sprintf("%d year ago", int(years))
	}
	return fmt.Sprintf("%d years ago", int(years))
}

func (s *Server) isAuthenticated(r *http.Request) bool {
	return s.sessionManager.Exists(r.Context(), "authenticatedUserID")
}
