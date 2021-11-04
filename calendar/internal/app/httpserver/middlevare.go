package httpserver

import (
	"calendar/internal/app"
	"calendar/internal/auth"
	"context"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func init() {
	if err := prometheus.Register(totalRequest); err != nil {
		log.Errorln("failed registrant request counter struct\n", err)
	}

	if err := prometheus.Register(responseStatus); err != nil {
		log.Errorln("failed registrant request status struct\n", err)
	}

	if err := prometheus.Register(httpDuration); err != nil {
		log.Errorln("failed registrant request duration struct\n", err)
	}
}

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func newResponseWriter(r http.ResponseWriter) *responseWriter {
	return &responseWriter{
		ResponseWriter: r,
		statusCode:     http.StatusOK,
	}
}

func (w responseWriter) WriteHeader(code int) {
	w.statusCode = code
	w.ResponseWriter.WriteHeader(w.statusCode)
}

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

var totalRequest = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "http_request_total",
		Help: "Number of get request",
	},
	[]string{"path"},
)

var responseStatus = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "response_status",
		Help: "Status of HTTP response",
	},
	[]string{"status"},
)

var httpDuration = prometheus.NewHistogramVec(
	prometheus.HistogramOpts{
		Name: "http_response_time_seconds",
		Help: "Duration of HTTP requests.",
	},
	[]string{"path"},
)

func (s Server) prometheusMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		route := mux.CurrentRoute(r)
		path, _ := route.GetPathTemplate()

		timer := prometheus.NewTimer(httpDuration.WithLabelValues(path))
		rw := newResponseWriter(w)
		next.ServeHTTP(rw, r)

		responseStatus.WithLabelValues(strconv.Itoa(rw.statusCode)).Inc()
		totalRequest.WithLabelValues(path).Inc()

		timer.ObserveDuration()
	})
}
