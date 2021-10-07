package httpserver

import (
	"calendar/internal/app/httpserver/auth"
	"calendar/internal/store"
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"net/http"
	"time"
)

type Server struct {
	store            store.Store
	router           *mux.Router
	logger           *logrus.Logger
	jwtWrapper       *auth.JwtWrapper
	userLocalization time.Location
}

const (
	jsonContentType     = "application/json"
	defaultLocalization = "Europe/Kiev"
)

func NewServer(store store.Store, wrapper *auth.JwtWrapper) *Server {
	s := &Server{
		store:      store,
		router:     mux.NewRouter(),
		logger:     logrus.New(),
		jwtWrapper: wrapper,
	}

	configureRouter(s)

	return s
}

func (s Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func (s Server) sendError(w http.ResponseWriter, statusCode int, data interface{}) {
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}

func (s Server) sendSuccess(w http.ResponseWriter, data interface{}) {
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(data)
}
