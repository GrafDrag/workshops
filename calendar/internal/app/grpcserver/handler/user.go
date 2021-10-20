package handler

import (
	"calendar/internal/app"
	"calendar/internal/auth"
	"calendar/internal/controller"
	"calendar/pb"
	"context"
)

type UserHandler struct {
	Server *app.IServer
}

func (u *UserHandler) Logout(ctx context.Context, request *pb.LogoutRequest) (*pb.LogoutResponse, error) {
	userID := ctx.Value(auth.KeyUserID).(int)
	userSession, err := u.Server.GetUserSession(userID)
	if err != nil {
		return nil, err
	}

	delete(userSession, u.Server.AuthToken)
	if err := u.Server.SetUserSession(userID, userSession); err != nil {
		return nil, err
	}

	return &pb.LogoutResponse{
		Status: pb.LogoutResponse_Successful,
	}, nil
}

func (u *UserHandler) Update(ctx context.Context, request *pb.UserUpdateRequest) (*pb.UserUpdateResponse, error) {
	c := controller.NewUserController(u.Server.Store)
	form := &controller.UpdateUserForm{
		Login:    request.GetLogin(),
		Timezone: request.GetTimezone(),
	}

	if err := c.Update(ctx, form); err != nil {
		return nil, err
	}

	return &pb.UserUpdateResponse{
		Status: pb.UserUpdateResponse_Successful,
	}, nil
}
