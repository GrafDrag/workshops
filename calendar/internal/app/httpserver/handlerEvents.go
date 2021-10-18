package httpserver

import (
	"calendar/internal/model"
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

func (s *Server) HandleListEvents(w http.ResponseWriter, r *http.Request) {
	user, err := s.store.User().FindById(r.Context().Value(KeyUserID).(int))
	if err != nil {
		s.sendError(w, http.StatusNotFound, errUserNotFound)
		return
	}

	q := r.URL.Query()
	searchModel := model.SearchEvent{
		UserID:   user.ID,
		Title:    q.Get("title"),
		Timezone: q.Get("timezone"),
		DateFrom: q.Get("dateFrom"),
		DateTo:   q.Get("dateTo"),
		TimeFrom: q.Get("timeFrom"),
		TimeTo:   q.Get("timeTo"),
	}

	res, err := s.store.Event().FindByParams(searchModel)

	if err != nil {
		s.sendError(w, http.StatusInternalServerError, err.Error())
	}

	s.sendSuccess(w, res)
}

func (s *Server) HandleGetEventsById(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		s.sendError(w, http.StatusBadRequest, errInvalidParams)
		return
	}

	event, err := s.store.Event().FindById(id)
	if err != nil || event.UserID != r.Context().Value(KeyUserID).(int) {
		s.sendError(w, http.StatusNotFound, errEventNotFound)
		return
	}

	s.sendSuccess(w, event)
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

	s.sendSuccessfullySaved(w)
}

func (s *Server) HandleUpdateEvent(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		s.sendError(w, http.StatusBadRequest, errInvalidParams)
		return
	}

	event, err := s.store.Event().FindById(id)
	if err != nil || event.UserID != r.Context().Value(KeyUserID).(int) {
		s.sendError(w, http.StatusNotFound, errEventNotFound)
		return
	}

	if err := json.NewDecoder(r.Body).Decode(event); err != nil {
		s.sendError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if err = s.store.Event().Update(event); err != nil {
		s.sendError(w, http.StatusInternalServerError, err.Error())
		return
	}

	s.sendSuccessfullySaved(w)
}

func (s *Server) HandleDeleteEvent(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		s.sendError(w, http.StatusBadRequest, errInvalidParams)
		return
	}

	event, err := s.store.Event().FindById(id)
	if err != nil || event.UserID != r.Context().Value(KeyUserID).(int) {
		s.sendError(w, http.StatusNotFound, errEventNotFound)
		return
	}

	if err := s.store.Event().Delete(event.ID); err != nil {
		s.sendError(w, http.StatusNotFound, errEventNotFound)
		return
	}

	s.sendSuccess(w, nil)
}
