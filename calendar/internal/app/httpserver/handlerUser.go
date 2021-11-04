package httpserver

import (
	"calendar/internal/auth"
	"calendar/internal/controller"
	"encoding/json"
	"net/http"
)

func (s *Server) HandleAuth(w http.ResponseWriter, r *http.Request) {
	form := &controller.LoginForm{}

	if err := json.NewDecoder(r.Body).Decode(form); err != nil {
		s.sendError(w, http.StatusBadRequest, err.Error())
		return
	}

	c := controller.NewAuthController(s.Store, &s.IServer)
	if token, err := c.Login(form); err != nil {
		s.sendError(w, err.GetStatus(), err.Error())
		return
	} else {
		form.Token = token
	}

	s.sendSuccess(w, form)
}

func (s *Server) HandleLogout(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(auth.KeyUserID).(int)
	userSession, err := s.GetUserSession(userID)
	if err != nil {
		s.sendError(w, http.StatusInternalServerError, err.Error())
		return
	}

	delete(userSession, s.AuthToken)
	if err := s.SetUserSession(userID, userSession); err != nil {
		s.sendError(w, http.StatusInternalServerError, err.Error())
		return
	}
}

func (s *Server) HandelUpdateUser(w http.ResponseWriter, r *http.Request) {
	form := &controller.UpdateUserForm{}

	if err := json.NewDecoder(r.Body).Decode(form); err != nil {
		s.sendError(w, http.StatusBadRequest, err.Error())
		return
	}

	c := controller.NewUserController(s.Store)

	if err := c.Update(r.Context(), form); err != nil {
		s.sendError(w, err.GetStatus(), err.Error())
		return
	}

	s.sendSuccessfullySaved(w)
}
