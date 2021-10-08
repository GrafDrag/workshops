package httpserver

import (
	"calendar/internal/app/httpserver/auth"
	"calendar/internal/store/inmemory"
	"github.com/sirupsen/logrus"
	"net/http"
)

var server *Server

func Start(config *Config) error {
	store := inmemory.New()
	wrapper := &auth.JwtWrapper{
		SecretKey:       config.JwtSecretKey,
		ExpirationHours: config.JwtExpHours,
		Issuer:          "AuthService",
	}
	server = NewServer(store, wrapper)

	logrus.Infof("Server starting on %v...", config.BindAddr)
	return http.ListenAndServe(config.BindAddr, server)
}