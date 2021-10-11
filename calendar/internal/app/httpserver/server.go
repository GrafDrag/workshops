package httpserver

import (
	"calendar/internal/app/httpserver/auth"
	"calendar/internal/store"
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"net/http"
)

type Server struct {
	store      store.Store
	router     *mux.Router
	logger     *logrus.Logger
	jwtWrapper *auth.JwtWrapper
}

const (
	JsonContentType = "application/json"
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

func (s Server) sendError(w http.ResponseWriter, statusCode int, errStr string) {
	w.WriteHeader(statusCode)
	err := json.NewEncoder(w).Encode(map[string]string{"error": errStr})
	if err != nil {
		s.logger.Error("failed encode response", err)
	}
}

func (s Server) sendSuccess(w http.ResponseWriter, data interface{}) {
	w.WriteHeader(http.StatusOK)
	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		s.logger.Error("failed encode response", err)
	}
}

func (s Server) sendSuccessfullySaved(w http.ResponseWriter) {
	w.WriteHeader(http.StatusCreated)
}
