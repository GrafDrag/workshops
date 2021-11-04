package httpserver

import (
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
	"net/http/pprof"
)

func configureRouter(s *Server) {
	s.router.Handle("/prometheus", promhttp.Handler())

	s.router.HandleFunc("/debug/pprof/", pprof.Index)
	s.router.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	s.router.HandleFunc("/debug/pprof/profile", pprof.Profile)
	s.router.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	s.router.HandleFunc("/debug/pprof/trace", pprof.Trace)

	s.router.Use(s.setContentType)
	s.router.Use(s.logRequest)

	s.router.Use(s.prometheusMiddleware)

	s.router.HandleFunc("/login", s.HandleAuth).Methods(http.MethodPost).Name("Login")
	s.router.Handle("/logout", s.authenticateUser(http.HandlerFunc(s.HandleLogout))).Methods(http.MethodGet).Name("Logout")

	s.router.HandleFunc(`/debug/pprof/{name:[\w]+}`, pprof.Index)

	auth := s.router.PathPrefix("/api").Subrouter()

	auth.Use(s.authenticateUser)
	auth.HandleFunc("/user", s.HandelUpdateUser).Methods(http.MethodPut).Name("Update user")
	auth.HandleFunc("/events", s.HandleListEvents).Methods(http.MethodGet).Name("Get list events")
	auth.HandleFunc("/events", s.HandleCreateEvent).Methods(http.MethodPost).Name("Create event")
	auth.HandleFunc(`/event/{id:[\d]+}`, s.HandleGetEventsById).Methods(http.MethodGet).Name("Get event by id")
	auth.HandleFunc("/event/{id:[0-9]+}", s.HandleUpdateEvent).Methods(http.MethodPut).Name("Update event")
	auth.HandleFunc("/event/{id:[0-9]+}", s.HandleDeleteEvent).Methods(http.MethodDelete).Name("Delete event")
}
