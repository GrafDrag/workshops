package httpserver

import (
	"fmt"
	"net/http"
)

const (
	methodGet    = "GET"
	methodPost   = "POST"
	methodPut    = "PUT"
	methodDelete = "DELETE"
)

func configureRouter(s *Server) {
	s.router.Use(s.setContentType)
	s.router.Use(s.logRequest)

	s.router.HandleFunc("/", Index).Methods(methodGet).Name("Index")
	s.router.HandleFunc("/login", s.HandleAuth).Methods(methodPost).Name("Login")
	s.router.HandleFunc("/logout", s.HandleLogout).Methods(methodGet).Name("Logout")

	auth := s.router.PathPrefix("/api").Subrouter()
	auth.Use(s.authenticateUser)

	auth.HandleFunc("/user", s.HandelUpdateUser).Methods(methodPut).Name("Update user")
	auth.HandleFunc("/events", s.HandleListEvents).Methods(methodGet).Name("Get list events")
	auth.HandleFunc("/events", s.HandleCreateEvent).Methods(methodPost).Name("Create event")
	auth.HandleFunc("/event/{id}", s.HandleGetEventsById).Methods(methodGet).Name("Get event by id")
	auth.HandleFunc("/event/{id}", s.HandleUpdateEvent).Methods(methodPut).Name("Update event")
	auth.HandleFunc("/event/{id}", s.HandleDeleteEvent).Methods(methodDelete).Name("Delete event")

}

func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello World!")
}
