package grpcserver

import (
	"calendar/internal/auth"
	"calendar/internal/config"
	"calendar/internal/session/redis"
	"calendar/internal/store/sqlstore"
	"github.com/sirupsen/logrus"
	"net"
)

var server *Server

func Start(config *config.Config) error {
	lis, err := net.Listen("tcp", config.GRPC.BindAddr)
	if err != nil {
		return err
	}

	store, close, err := sqlstore.New(config.DB)
	if err != nil {
		return err
	}
	defer close()

	wrapper := auth.NewJwtWrapper(config.Jwt, "AuthService")

	session, err := redis.NewSession(config.Session)
	if err != nil {
		return err
	}
	server = NewServer(store, wrapper, session)

	logrus.Infof("Server starting on %v...", config.GRPC.BindAddr)

	return server.Serve(lis)
}
