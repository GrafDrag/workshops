package httpserver

import (
	"fmt"
	"net/http"
)

func configureRouter(s *Server) {
	s.router.Use(s.setContentType)
	s.router.Use(s.logRequest)

	s.router.HandleFunc("/", Index).Methods(http.MethodGet).Name("Index")
	s.router.HandleFunc("/login", s.HandleAuth).Methods(http.MethodPost).Name("Login")
	s.router.HandleFunc("/logout", s.HandleLogout).Methods(http.MethodGet).Name("Logout")

	auth := s.router.PathPrefix("/api").Subrouter()
	auth.Use(s.authenticateUser)

	auth.HandleFunc("/user", s.HandelUpdateUser).Methods(http.MethodPut).Name("Update user")
	auth.HandleFunc("/events", s.HandleListEvents).Methods(http.MethodGet).Name("Get list events")
	auth.HandleFunc("/events", s.HandleCreateEvent).Methods(http.MethodPost).Name("Create event")
	auth.HandleFunc("/event/{id}", s.HandleGetEventsById).Methods(http.MethodGet).Name("Get event by id")
	auth.HandleFunc("/event/{id}", s.HandleUpdateEvent).Methods(http.MethodPut).Name("Update event")
	auth.HandleFunc("/event/{id}", s.HandleDeleteEvent).Methods(http.MethodDelete).Name("Delete event")
}

func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello World!")
}
