package httpserver

import (
	"calendar/internal/auth"
	"calendar/internal/config"
	"calendar/internal/session/redis"
	"calendar/internal/store/sqlstore"
	"github.com/sirupsen/logrus"
	"net/http"
)

var server *Server

func Start(config *config.Config) error {
	store, close, err := sqlstore.New(config.DB)
	if err != nil {
		return err
	}
	defer func() {
		if err := close(); err != nil {
			return
		}
	}()

	wrapper := auth.NewJwtWrapper(config.Jwt, "AuthService")

	session, err := redis.NewSession(config.Session)
	if err != nil {
		return err
	}

	server = NewServer(store, session, wrapper)

	logrus.Infof("Server starting on %v...", config.Rest.BindAddr)
	return http.ListenAndServe(config.Rest.BindAddr, server)
}
