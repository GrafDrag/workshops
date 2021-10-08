package httpserver

import (
	"calendar/internal/model"
	"encoding/json"
	"fmt"
	"net/http"
)

func (s *Server) HandleListEvents(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "HandleList")
}

func (s *Server) HandleGetEventsById(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "HandleGet")
}

func (s *Server) HandleCreateEvent(w http.ResponseWriter, r *http.Request) {
	event := &model.Event{}
	if err := json.NewDecoder(r.Body).Decode(event); err != nil {
		s.sendError(w, http.StatusInternalServerError, err.Error())
	}

	user, err := s.store.User().FindById(r.Context().Value(KeyUserID).(int))
	if err != nil {
		s.sendError(w, http.StatusNotFound, errUserNotFound)
		return
	}

	event.UserID = user.ID
	if err := event.Validate(); err != nil {
		s.sendError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := s.store.Event().Create(event); err != nil {
		s.sendError(w, http.StatusInternalServerError, err.Error())
		return
	}

	s.sendSuccess(w, nil)
}

func (s *Server) HandleUpdateEvent(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "HandleUpdate")
}

func (s *Server) HandleDeleteEvent(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "HandleDeleteEvent")
}
