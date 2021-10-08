package httpserver_test

import (
	"bytes"
	"calendar/internal/app/httpserver"
	"calendar/internal/app/httpserver/auth"
	"calendar/internal/model"
	"calendar/internal/store/inmemory"
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

var (
	store      = inmemory.New()
	jwtWrapper = &auth.JwtWrapper{
		SecretKey:       "secretkey",
		Issuer:          "AuthService",
		ExpirationHours: 24,
	}
)

func TestUser_Login(t *testing.T) {
	server := mustMakeServer()

	t.Run("test user login", func(t *testing.T) {
		user := model.TestUser(t)
		rec := httptest.NewRecorder()
		request := newGetLoginRequest(t, user.Login, user.Password)
		server.ServeHTTP(rec, request)

		assert.Equal(t, rec.Code, http.StatusOK)
		assert.Equal(t, rec.Header().Get("content-type"), httpserver.JsonContentType)
	})
}

func TestUser_Logout(t *testing.T) {
	server := mustMakeServer()

	t.Run("test user logout", func(t *testing.T) {
		user := model.TestUser(t)
		token := getJWTToken(t, user)
		rec := httptest.NewRecorder()

		request, _ := http.NewRequest(http.MethodGet, "/logout", nil)
		request.Header.Set("Authorization", token)

		server.ServeHTTP(rec, request)

		assert.Equal(t, rec.Code, http.StatusOK)
		assert.Equal(t, rec.Header().Get("content-type"), httpserver.JsonContentType)
	})
}

func TestUser_Update(t *testing.T) {
	server := mustMakeServer()

	t.Run("update user timezone", func(t *testing.T) {
		timezone := "Europe/Kiev"
		user := model.TestUser(t)
		token := getJWTToken(nil, user)
		rec := httptest.NewRecorder()

		b := &bytes.Buffer{}
		err := json.NewEncoder(b).Encode(map[string]interface{}{
			"login":    user.Login,
			"timezone": timezone,
		})
		if err != nil {
			t.Fatalf("could not encode data. %s", err.Error())
		}

		request, _ := http.NewRequest(http.MethodPut, "/api/user", b)
		request.Header.Set("Authorization", fmt.Sprintf("Bearer %v", token))

		server.ServeHTTP(rec, request)

		assert.Equal(t, rec.Code, http.StatusOK)
		assert.Equal(t, rec.Header().Get("content-type"), httpserver.JsonContentType)

		got, _ := store.User().FindById(user.ID)
		assert.Equal(t, timezone, got.Timezone)
	})

	t.Run("update user login", func(t *testing.T) {
		newLogin := "testUser2"
		user := model.TestUser(t)
		token := getJWTToken(nil, user)
		rec := httptest.NewRecorder()

		b := &bytes.Buffer{}
		err := json.NewEncoder(b).Encode(map[string]interface{}{
			"login":    newLogin,
			"timezone": user.Timezone,
		})
		if err != nil {
			t.Fatalf("could not encode data. %s", err.Error())
		}

		request, _ := http.NewRequest(http.MethodPut, "/api/user", b)
		request.Header.Set("Authorization", fmt.Sprintf("Bearer %v", token))

		server.ServeHTTP(rec, request)

		assert.Equal(t, rec.Code, http.StatusOK)
		assert.Equal(t, rec.Header().Get("content-type"), httpserver.JsonContentType)

		got, _ := store.User().FindById(user.ID)
		assert.Equal(t, newLogin, got.Login)
	})
}

func TestServer_HandleListEvents(t *testing.T) {

}

func TestServer_HandleCreateEvent(t *testing.T) {

}

func TestServer_HandleGetEventsById(t *testing.T) {

}

func TestServer_HandleUpdateEvent(t *testing.T) {

}

func TestServer_HandleDeleteEvent(t *testing.T) {

}

func mustMakeServer() *httpserver.Server {
	return httpserver.NewServer(store, jwtWrapper)
}

func newGetLoginRequest(t *testing.T, login, password string) *http.Request {
	b := &bytes.Buffer{}
	err := json.NewEncoder(b).Encode(map[string]interface{}{
		"login":    login,
		"password": password,
	})
	if err != nil {
		t.Fatalf("could not encode data. %s", err.Error())
	}

	req, _ := http.NewRequest(http.MethodPost, "/login", b)

	return req
}

func getJWTToken(t *testing.T, user *model.User) string {
	if err := store.User().Create(user); err != nil {
		t.Fatal("could not create user")
	}
	token, _ := jwtWrapper.GenerateToken(user)

	return token
}
