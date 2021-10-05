package handler

import (
	"fmt"
	"net/http"
)

func HandleListEvents(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "HandleList")
}

func HandleGetEventsById(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "HandleGet")
}

func HandleCreateEvent(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "HandleCreate")
}

func HandleUpdateEvent(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "HandleUpdate")
}

func HandleDeleteEvent(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "HandleDeleteEvent")
}
