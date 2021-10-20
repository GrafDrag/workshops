package controller

import (
	"calendar/internal/auth"
	"calendar/internal/model"
	"calendar/internal/store"
	"context"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"net/http"
)

type UserController struct {
	store store.Store
}

func NewUserController(store store.Store) *UserController {
	return &UserController{
		store: store,
	}
}

type UpdateUserForm struct {
	Login    string `json:"login"`
	Timezone string `json:"timezone"`
}

func (f UpdateUserForm) Validate() error {
	return validation.ValidateStruct(
		&f,
		validation.Field(&f.Login, validation.When(f.Login != "", validation.Length(6, 20))),
		validation.Field(&f.Timezone, validation.By(model.TimeZoneValidator(f.Timezone))),
	)
}

func (c *UserController) Update(ctx context.Context, form *UpdateUserForm) *ResponseError {
	repository := c.store.User()
	if err := form.Validate(); err != nil {
		return &ResponseError{
			Err: err,
		}
	}

	userID := ctx.Value(auth.KeyUserID).(int)
	u, err := repository.FindById(userID)
	if err != nil {
		return &ResponseError{
			status: http.StatusBadRequest,
			Err:    errUserNotFound,
		}
	}

	if form.Login != "" {
		if form.Login != u.Login {
			if _, err := repository.FindByLogin(form.Login); err == nil {
				return &ResponseError{
					status: http.StatusBadRequest,
					Err:    errUserExist,
				}
			}
		}

		u.Login = form.Login
	}

	if form.Timezone != "" {
		u.Timezone = form.Timezone
	}

	if err := repository.Update(u); err != nil {
		return &ResponseError{
			Err: err,
		}
	}
	return nil
}
