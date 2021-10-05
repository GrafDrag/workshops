package httpserver

import (
	"github.com/sirupsen/logrus"
	"net/http"
	"time"
)

func Logger(handler http.HandlerFunc, name string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		handler.ServeHTTP(w, r)
		logrus.Infof("method: %s uri: %s  name: %s (%v)", r.Method, r.RequestURI, name, time.Now())
	}
}
