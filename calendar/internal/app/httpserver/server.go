package httpserver

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"net/http"
	"time"
)

type Server struct {
	router *mux.Router
	logger *logrus.Logger
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

func (s Server) setLocalization(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		loc, err := time.LoadLocation(defaultLocalization)
		if err != nil {
			s.logger.Fatal(err)
		}

		fmt.Println(time.Now().In(loc))
		next.ServeHTTP(w, r)
	})
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
