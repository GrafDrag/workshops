package httpserver

import (
	"context"
	"net/http"
	"strings"
	"time"
)

type ctxUserKey int

const (
	KeyUserID ctxUserKey = iota
)

func (s *Server) authenticateUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenHeader := r.Header.Get("Authorization")

		if tokenHeader == "" {
			s.sendError(w, http.StatusForbidden, errEmptyAuthToken)
			return
		}

		splitted := strings.Split(tokenHeader, " ")
		if len(splitted) != 2 {
			s.sendError(w, http.StatusForbidden, errInvalidAuthToken)
			return
		}

		claims, err := s.jwtWrapper.ValidateToken(splitted[1])
		if err != nil {
			s.sendError(w, http.StatusForbidden, errInvalidAuthToken)
			return
		}

		s.logger.Infof("User auth by token: %v", claims.ID)
		ctx := context.WithValue(r.Context(), KeyUserID, claims.ID)
		r = r.WithContext(ctx)

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
		s.logger.Infof("started %s %s", r.Method, r.RequestURI)

		start := time.Now()

		next.ServeHTTP(w, r)

		s.logger.Infof("completed with in %v", time.Since(start))
	})
}
