package handler

import (
	"calendar/internal/app"
	"calendar/internal/controller"
	"calendar/pb"
	"context"
)

type AuthHandler struct {
	Server *app.IServer
}

func (s *AuthHandler) Login(ctx context.Context, in *pb.LoginRequest) (*pb.LoginResponse, error) {
	c := controller.NewAuthController(s.Server.Store, s.Server)
	form := &controller.LoginForm{
		Login:    in.GetLogin(),
		Password: in.GetPassword(),
	}

	token, err := c.Login(form)
	if err != nil {
		return nil, err
	}

	return &pb.LoginResponse{
		Token: token,
	}, nil
}
