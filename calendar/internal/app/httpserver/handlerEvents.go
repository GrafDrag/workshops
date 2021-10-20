package httpserver

import (
	"calendar/internal/app"
	"calendar/internal/controller"
	"calendar/internal/model"
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
)

func (s *Server) HandleListEvents(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	searchModel := model.SearchEvent{
		Title:    q.Get("title"),
		Timezone: q.Get("timezone"),
		DateFrom: q.Get("dateFrom"),
		DateTo:   q.Get("dateTo"),
		TimeFrom: q.Get("timeFrom"),
		TimeTo:   q.Get("timeTo"),
	}

	c := controller.NewEventController(s.Store)
	res, err := c.List(r.Context(), searchModel)

	if err != nil {
		s.sendError(w, err.GetStatus(), err.Error())
		return
	}

	s.sendSuccess(w, res)
}

func (s *Server) HandleGetEventsById(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	c := controller.NewEventController(s.Store)

	event, err := c.FindById(r.Context(), params["id"])
	if err != nil {
		s.sendError(w, err.GetStatus(), err.Error())
		return
	}

	s.sendSuccess(w, event)
}

func (s *Server) HandleCreateEvent(w http.ResponseWriter, r *http.Request) {
	event := &model.Event{}
	if err := json.NewDecoder(r.Body).Decode(event); err != nil {
		s.sendError(w, http.StatusInternalServerError, err.Error())
	}

	c := controller.NewEventController(s.Store)

	if err := c.Create(r.Context(), event); err != nil {
		s.sendError(w, err.GetStatus(), err.Error())
		return
	}

	s.sendSuccessfullySaved(w)
}

func (s *Server) HandleUpdateEvent(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	c := controller.NewEventController(s.Store)

	event, err := c.FindById(r.Context(), params["id"])
	if err != nil {
		s.sendError(w, err.GetStatus(), err.Error())
	}

	if err := json.NewDecoder(r.Body).Decode(event); err != nil {
		s.sendError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if err := c.Update(r.Context(), event); err != nil {
		s.sendError(w, err.GetStatus(), err.Error())
		return
	}

	s.sendSuccessfullySaved(w)
}

func (s *Server) HandleDeleteEvent(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	c := controller.NewEventController(s.Store)

	if err := c.Delete(r.Context(), params["id"]); err != nil {
		s.sendError(w, http.StatusNotFound, app.errEventNotFound)
		return
	}

	s.sendSuccess(w, nil)
}
