package httpserver

import (
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
	fmt.Fprintf(w, "HandleCreate")
}

func (s *Server) HandleUpdateEvent(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "HandleUpdate")
}

func (s *Server) HandleDeleteEvent(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "HandleDeleteEvent")
}
