package httpserver

import (
	"calendar/internal/model"
	"encoding/json"
	"fmt"
	validation "github.com/go-ozzo/ozzo-validation"
	"net/http"
)

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

func (s *Server) HandleAuth(w http.ResponseWriter, r *http.Request) {
	form := &LoginForm{}

	if err := json.NewDecoder(r.Body).Decode(form); err != nil {
		s.sendError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := form.Validate(); err != nil {
		s.sendError(w, http.StatusBadRequest, err.Error())
		return
	}

	u, err := s.store.User().FindByLogin(form.Login)
	if err != nil {
		u = &model.User{
			Login:    form.Login,
			Password: form.Password,
		}

		if err = s.store.User().Create(u); err != nil {
			s.sendError(w, http.StatusInternalServerError, err.Error())
			return
		}
	}

	if err = u.CheckPassword(form.Password); err != nil {
		s.sendError(w, http.StatusBadRequest, errIncorrectLoginOrPassword)
		return
	}

	form.Password = ""
	form.Token, err = s.jwtWrapper.GenerateToken(u)
	if err != nil {
		s.sendError(w, http.StatusInternalServerError, err.Error())
		return
	}

	s.sendSuccess(w, form)
}

func (s *Server) HandleLogout(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "HandleLogout")
}

type UpdateUserForm struct {
	Login    string `json:"login"`
	Timezone string `json:"timezone"`
}

func (f UpdateUserForm) Validate() error {
	return validation.ValidateStruct(
		&f,
		validation.Field(&f.Login, validation.Length(6, 20)),
		validation.Field(&f.Timezone, validation.By(model.TimeZoneValidator(f.Timezone))),
	)
}

func (s *Server) HandelUpdateUser(w http.ResponseWriter, r *http.Request) {
	form := &UpdateUserForm{}

	if err := json.NewDecoder(r.Body).Decode(form); err != nil {
		s.sendError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := form.Validate(); err != nil {
		s.sendError(w, http.StatusBadRequest, err.Error())
		return
	}

	userID := r.Context().Value(KeyUserID).(int)
	u, err := s.store.User().FindById(userID)
	if err != nil {
		s.sendError(w, http.StatusNotFound, errUserNotFound)
		return
	}

	if form.Login != u.Login {
		if _, err := s.store.User().FindByLogin(form.Login); err == nil {
			s.sendError(w, http.StatusBadRequest, errUserExist)
			return
		}
	}

	u.Login = form.Login
	u.Timezone = form.Timezone

	if err := s.store.User().Update(u); err != nil {
		s.sendError(w, http.StatusInternalServerError, err.Error())
		return
	}

	s.sendSuccess(w, form)
}
