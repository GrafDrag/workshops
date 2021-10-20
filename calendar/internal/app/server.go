package app

import (
	"calendar/internal/auth"
	"calendar/internal/session"
	"calendar/internal/store"
	"encoding/json"
	"github.com/sirupsen/logrus"
	"strconv"
)

type IServer struct {
	Store      store.Store
	Session    session.Session
	Logger     *logrus.Logger
	JWTWrapper *auth.JwtWrapper
	AuthToken  string
}

func (s *IServer) GetUserSession(ID int) (map[string]bool, error) {
	var userSession map[string]bool
	v, _ := s.Session.Get(strconv.Itoa(ID))

	if v == "" {
		userSession = map[string]bool{}
	} else if err := json.Unmarshal([]byte(v), &userSession); err != nil {
		return nil, err
	}

	return userSession, nil
}

func (s *IServer) SetUserSession(ID int, userSession map[string]bool) error {
	jb, err := json.Marshal(userSession)

	if err != nil {
		return err
	}

	if err := s.Session.Set(strconv.Itoa(ID), string(jb)); err != nil {
		return err
	}

	return nil
}
