package grpcserver

import (
	"calendar/internal/app"
	"calendar/internal/app/grpcserver/handler"
	"calendar/internal/app/grpcserver/interceptor"
	"calendar/internal/auth"
	"calendar/internal/session"
	"calendar/internal/store"
	"calendar/pb"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

type Server struct {
	app.IServer
	*grpc.Server
}

func NewServer(store store.Store, wrapper *auth.JwtWrapper, session session.Session) *Server {
	s := &Server{}

	s.Store = store
	s.Session = session
	s.JWTWrapper = wrapper
	s.Logger = logrus.New()

	authInterceptor := interceptor.NewAuthInterceptor(&s.IServer)

	s.Server = grpc.NewServer(
		grpc.UnaryInterceptor(authInterceptor.Unary()),
		grpc.StreamInterceptor(authInterceptor.Stream()),
	)

	pb.RegisterEventServiceServer(s.Server, &handler.EventHandler{Server: &s.IServer})
	pb.RegisterUserServiceServer(s.Server, &handler.UserHandler{Server: &s.IServer})
	pb.RegisterAuthServiceServer(s.Server, &handler.AuthHandler{Server: &s.IServer})

	return s
}
