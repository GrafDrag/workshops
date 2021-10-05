package httpserver

import (
	"github.com/sirupsen/logrus"
	"net/http"
)

func Start(config *Config) error {
	s := NewServer()

	logrus.Infof("Server starting on %v...", config.BindAddr)
	return http.ListenAndServe(config.BindAddr, s)
}
