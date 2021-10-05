package handler

import (
	"fmt"
	"net/http"
)

func HandleAuth(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "HandleAuth")
}

func HandleLogout(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "HandleLogout")
}

func HandelUpdateUser(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "HandelUpdateUser")
}
