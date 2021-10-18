package httpserver

import (
	"calendar/internal/app/httpserver/auth"
	"calendar/internal/session"
	"calendar/internal/store"
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"
)

type Server struct {
	store      store.Store
	session    session.Session
	router     *mux.Router
	logger     *logrus.Logger
	jwtWrapper *auth.JwtWrapper
	authToken  string
}

const (
	JsonContentType = "application/json"
)

func NewServer(store store.Store, wrapper *auth.JwtWrapper, session session.Session) *Server {
	s := &Server{
		store:      store,
		session:    session,
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

func (s *Server) getUserSession(ID int) (map[string]bool, error) {
	var userSession map[string]bool
	v, _ := s.session.Get(strconv.Itoa(ID))

	if v == "" {
		userSession = map[string]bool{}
	} else if err := json.Unmarshal([]byte(v), &userSession); err != nil {
		return nil, err
	}

	return userSession, nil
}

func (s *Server) setUserSession(ID int, userSession map[string]bool) error {
	jb, err := json.Marshal(userSession)

	if err != nil {
		return err
	}

	if err := s.session.Set(strconv.Itoa(ID), string(jb)); err != nil {
		return err
	}

	return nil
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
