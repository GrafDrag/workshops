package httpserver

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"net/http"
	"time"
)

type Server struct {
	userLocalization time.Location
	router           *mux.Router
	logger           *logrus.Logger
}

const (
	jsonContentType     = "application/json"
	defaultLocalization = "America/Chicago"
)

func NewServer() *Server {
	s := &Server{
		router: mux.NewRouter(),
		logger: logrus.New(),
	}

	configureRouter(s)

	return s
}

func (s Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func (s *Server) authenticateUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("auth user")

		next.ServeHTTP(w, r)
	})
}

func (s Server) setContentType(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("content-type", jsonContentType)

		next.ServeHTTP(w, r)
	})
}

func (s Server) logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		s.logger.Infof("started %s %s", r.Method, r.RequestURI)

		start := time.Now()

		next.ServeHTTP(w, r)

		s.logger.Infof("completed with in %v", time.Now().Sub(start))
	})
}
