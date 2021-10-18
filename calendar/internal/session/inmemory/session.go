package inmemory

import (
	"errors"
	"fmt"
)

type session struct {
	sessions map[string]string
}

func NewSession() *session {
	return &session{
		sessions: map[string]string{},
	}
}

func (s *session) Get(key string) (string, error) {
	v, ok := s.sessions[key]
	if !ok {
		return "", errors.New(fmt.Sprintf("failed find key %s", key))
	}
	return v, nil
}

func (s *session) Set(key string, value string) error {
	s.sessions[key] = value
	return nil
}

func (s *session) Remove(key string) (string, error) {
	v, err := s.Get(key)
	if err != nil {
		return "", err
	}
	delete(s.sessions, key)
	return v, nil
}

func (s *session) Flash() error {
	s.sessions = map[string]string{}
	return nil
}
