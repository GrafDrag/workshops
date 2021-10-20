package httpserver

import (
	"calendar/internal/app"
	"calendar/internal/auth"
	"calendar/internal/session"
	"calendar/internal/store"
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"net/http"
)

type Server struct {
	app.IServer
	router *mux.Router
}

const (
	JsonContentType = "application/json"
)

func NewServer(store store.Store, session session.Session, wrapper *auth.JwtWrapper) *Server {
	s := &Server{
		router: mux.NewRouter(),
	}

	s.Store = store
	s.Session = session
	s.JWTWrapper = wrapper
	s.Logger = logrus.New()

	configureRouter(s)

	return s
}

func (s Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func (s Server) sendError(w http.ResponseWriter, status int, errorStr string) {
	w.WriteHeader(status)
	err := json.NewEncoder(w).Encode(map[string]string{"error": errorStr})
	if err != nil {
		s.Logger.Error("failed encode response", err)
	}
}

func (s Server) sendSuccess(w http.ResponseWriter, data interface{}) {
	w.WriteHeader(http.StatusOK)
	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		s.Logger.Error("failed encode response", err)
	}
}

func (s Server) sendSuccessfullySaved(w http.ResponseWriter) {
	w.WriteHeader(http.StatusCreated)
}
