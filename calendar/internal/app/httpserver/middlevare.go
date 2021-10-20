package httpserver

import (
	"calendar/internal/app"
	"calendar/internal/auth"
	"context"
	"net/http"
	"strings"
	"time"
)

func (s *Server) authenticateUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenHeader := r.Header.Get("Authorization")

		if tokenHeader == "" {
			s.sendError(w, http.StatusForbidden, app.ErrEmptyAuthToken)
			return
		}

		splitted := strings.Split(tokenHeader, " ")
		if len(splitted) != 2 {
			s.sendError(w, http.StatusForbidden, app.ErrInvalidAuthToken)
			return
		}

		claims, err := s.JWTWrapper.ValidateToken(splitted[1])
		if err != nil {
			s.sendError(w, http.StatusForbidden, app.ErrInvalidAuthToken)
			return
		}

		userSession, err := s.GetUserSession(claims.ID)
		if err != nil {
			s.sendError(w, http.StatusForbidden, app.ErrSessionNotFound)
			return
		}

		if _, ok := userSession[splitted[1]]; !ok {
			s.sendError(w, http.StatusForbidden, app.ErrInvalidAuthToken)
			return
		}

		s.Logger.Infof("User ID #%v auth by token", claims.ID)
		ctxUserID := context.WithValue(r.Context(), auth.KeyUserID, claims.ID)
		r = r.WithContext(ctxUserID)
		s.AuthToken = splitted[1]

		next.ServeHTTP(w, r)
	})
}

func (s Server) setContentType(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("content-type", JsonContentType)

		next.ServeHTTP(w, r)
	})
}

func (s Server) logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		s.Logger.Infof("started %s %s", r.Method, r.RequestURI)

		start := time.Now()

		next.ServeHTTP(w, r)

		s.Logger.Infof("completed with in %v", time.Since(start))
	})
}
