package controller

import (
	"calendar/internal/app"
	"calendar/internal/model"
	"calendar/internal/store"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"net/http"
)

type AuthController struct {
	store  store.Store
	server *app.IServer
}

func NewAuthController(store store.Store, server *app.IServer) *AuthController {
	return &AuthController{
		store:  store,
		server: server,
	}
}

type LoginForm struct {
	Login    string `json:"login"`
	Password string `json:"password"`
	Token    string `json:"token"`
}

func (f LoginForm) Validate() error {
	return validation.ValidateStruct(
		&f,
		validation.Field(&f.Login, validation.Required, validation.Length(6, 20)),
		validation.Field(&f.Password, validation.Required, validation.Length(6, 20)),
	)
}

func (c AuthController) Login(form *LoginForm) (string, *ResponseError) {
	repository := c.store.User()
	if err := form.Validate(); err != nil {
		return "", &ResponseError{
			status: http.StatusBadRequest,
			Err:    err,
		}
	}

	u, err := repository.FindByLogin(form.Login)
	if err != nil {
		u = &model.User{
			Login:    form.Login,
			Password: form.Password,
		}

		if err = repository.Create(u); err != nil {
			return "", &ResponseError{
				Err: err,
			}
		}
	}

	if err = u.CheckPassword(form.Password); err != nil {
		return "", &ResponseError{
			status: http.StatusBadRequest,
			Err:    errIncorrectLoginOrPassword,
		}
	}

	form.Password = ""
	token, err := c.server.JWTWrapper.GenerateToken(u)
	if err != nil {
		return "", &ResponseError{
			Err: err,
		}
	}

	userSession, err := c.server.GetUserSession(u.ID)
	if err != nil {
		return "", &ResponseError{
			Err: err,
		}
	}

	userSession[token] = true
	err = c.server.SetUserSession(u.ID, userSession)
	if err != nil {
		return "", &ResponseError{
			Err: err,
		}
	}

	return token, nil
}
