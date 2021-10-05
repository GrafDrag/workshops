package httpserver

import (
	"calendar/internal/app/httpserver/handler"
	"fmt"
	"net/http"
)

const (
	methodGet    = "GET"
	methodPost   = "POST"
	methodPut    = "PUT"
	methodDelete = "DELETE"
)

type Route struct {
	Name        string
	Method      string
	Path        string
	HandlerFunc http.HandlerFunc
}

type Routes []Route

func configureRouter(s *Server) {
	s.router.Use(s.setContentType)
	s.router.Use(s.logRequest)

	for _, route := range routes {
		s.router.HandleFunc(route.Path, route.HandlerFunc).Methods(route.Method).Name(route.Name)
	}

	auth := s.router.PathPrefix("/api").Subrouter()
	auth.Use(s.authenticateUser)
	for _, route := range authRoutes {
		auth.HandleFunc(route.Path, route.HandlerFunc).Methods(route.Method).Name(route.Name)
	}
}

func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello World!")
}

var routes = Routes{
	Route{
		"Index",
		methodGet,
		"/",
		Index,
	},

	Route{
		"Login",
		methodPost,
		"/login",
		handler.HandleAuth,
	},

	Route{
		"Logout",
		methodGet,
		"/logout",
		handler.HandleLogout,
	},
}

var authRoutes = Routes{
	Route{
		"Update user",
		methodPut,
		"/user",
		handler.HandelUpdateUser,
	},

	Route{
		"Get list events",
		methodGet,
		"/events",
		handler.HandleListEvents,
	},

	Route{
		"Create event",
		methodPost,
		"/events",
		handler.HandleCreateEvent,
	},

	Route{
		"Get event by id",
		methodGet,
		"/event/{id}",
		handler.HandleGetEventsById,
	},

	Route{
		"Update event",
		methodPut,
		"/event/{id}",
		handler.HandleUpdateEvent,
	},

	Route{
		"Delete event",
		methodDelete,
		"/event/{id}",
		handler.HandleDeleteEvent,
	},
}
