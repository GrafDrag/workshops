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
	"reflect"
	"strconv"
	"testing"
	"time"
)

var (
	store      = inmemory.New()
	jwtWrapper = &auth.JwtWrapper{
		SecretKey:       "secretkey",
		Issuer:          "AuthService",
		ExpirationHours: 24,
	}
	session = inmemory.NewSession()
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
		request.Header.Set("Authorization", fmt.Sprintf("Bearer %v", token))

		server.ServeHTTP(rec, request)

		assert.Equal(t, rec.Code, http.StatusOK)
		assert.Equal(t, rec.Header().Get("content-type"), httpserver.JsonContentType)

		server.ServeHTTP(rec, request)
		assert.Equal(t, rec.Code, http.StatusForbidden)
	})
}

func TestUser_Update(t *testing.T) {
	server := mustMakeServer()

	t.Run("update user timezone", func(t *testing.T) {
		timezone := "Europe/Kiev"
		user := model.TestUser(t)
		token := getJWTToken(t, user)
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

		assert.Equal(t, rec.Code, http.StatusCreated)
		assert.Equal(t, rec.Header().Get("content-type"), httpserver.JsonContentType)

		got, _ := store.User().FindById(user.ID)
		assert.Equal(t, timezone, got.Timezone)
	})

	t.Run("update user login", func(t *testing.T) {
		newLogin := "testUser2"
		user := model.TestUser(t)
		token := getJWTToken(t, user)
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

		assert.Equal(t, rec.Code, http.StatusCreated)
		assert.Equal(t, rec.Header().Get("content-type"), httpserver.JsonContentType)

		got, _ := store.User().FindById(user.ID)
		assert.Equal(t, newLogin, got.Login)
	})
}

func TestServer_HandleListEvents(t *testing.T) {
	server := mustMakeServer()
	user := model.TestUser(t)
	token := getJWTToken(t, user)
	event := model.TestEvent(t)
	event.UserID = user.ID
	event.Timezone = "Europe/Kiev"
	addEventToStore(t, event)
	var cases = []struct {
		name        string
		searchQuery func() string
		isEmpty     bool
	}{
		{
			name: "search by title",
			searchQuery: func() string {
				return fmt.Sprintf("%s=%s", "title", event.Title)
			},
			isEmpty: false,
		},
		{
			name: "search by timezone",
			searchQuery: func() string {
				return fmt.Sprintf("%s=%s", "timezone", event.Timezone)
			},
			isEmpty: false,
		},
		{
			name: "search by dateFrom",
			searchQuery: func() string {
				return fmt.Sprintf("%s=%s", "dateFrom", time.Now().Format("2006-01-02"))
			},
			isEmpty: false,
		},
		{
			name: "search by dateTo",
			searchQuery: func() string {
				return fmt.Sprintf("%s=%s", "dateTo", time.Now().Format("2006-01-02"))
			},
			isEmpty: true,
		},
		{
			name: "search by dateFrom with timeFrom",
			searchQuery: func() string {
				now := time.Now().Add(time.Minute * -10)
				return fmt.Sprintf(
					"%s=%s&%s=%s",
					"dateFrom",
					now.Format("2006-01-02"),
					"timeFrom",
					now.Format("15:04"),
				)
			},
			isEmpty: false,
		},
		{
			name: "search by dateTo with timeTo",
			searchQuery: func() string {
				now := time.Now().Add(time.Minute * 20)
				return fmt.Sprintf(
					"%s=%s&%s=%s",
					"dateTo",
					now.Format("2006-01-02"),
					"timeTo",
					now.Format("15:04"),
				)
			},
			isEmpty: false,
		},
	}
	for _, s := range cases {
		t.Run(s.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			request, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/api/events?%s", s.searchQuery()), nil)
			request.Header.Set("Authorization", fmt.Sprintf("Bearer %v", token))

			server.ServeHTTP(rec, request)

			assert.Equal(t, rec.Code, http.StatusOK)
			assert.Equal(t, rec.Header().Get("content-type"), httpserver.JsonContentType)

			var res []model.Event
			err := json.NewDecoder(rec.Body).Decode(&res)
			if err != nil {
				t.Fatalf("failed parse json %v", err)
			}

			assert.Equal(t, s.isEmpty, len(res) == 0)
		})
	}
}

func TestServer_HandleCreateEvent(t *testing.T) {
	server := mustMakeServer()

	t.Run("test event create", func(t *testing.T) {
		event := model.TestEvent(t)
		rec := httptest.NewRecorder()
		token := getJWTToken(t, model.TestUser(t))

		b := &bytes.Buffer{}
		err := json.NewEncoder(b).Encode(event)
		if err != nil {
			t.Fatalf("could not encode data. %s", err.Error())
		}

		request, _ := http.NewRequest(http.MethodPost, "/api/events", b)
		request.Header.Set("Authorization", fmt.Sprintf("Bearer %v", token))

		server.ServeHTTP(rec, request)

		assert.Equal(t, rec.Code, http.StatusCreated)
		assert.Equal(t, rec.Header().Get("content-type"), httpserver.JsonContentType)
	})
}

func TestServer_HandleGetEventsById(t *testing.T) {
	server := mustMakeServer()

	t.Run("test get event by id", func(t *testing.T) {
		event := model.TestEvent(t)
		user := model.TestUser(t)
		rec := httptest.NewRecorder()

		token := getJWTToken(t, user)
		event.UserID = user.ID

		addEventToStore(t, event)

		request, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/api/event/%d", event.ID), nil)
		request.Header.Set("Authorization", fmt.Sprintf("Bearer %v", token))

		server.ServeHTTP(rec, request)

		assert.Equal(t, rec.Code, http.StatusOK)
		assert.Equal(t, rec.Header().Get("content-type"), httpserver.JsonContentType)

		resp := &model.Event{}

		if err := json.NewDecoder(rec.Body).Decode(resp); err != nil {
			t.Fatalf("could not encode data. %s", err.Error())
		}

		assert.Equal(t, event.ID, resp.ID)
	})
}

func TestServer_HandleUpdateEvent(t *testing.T) {

	cases := []struct {
		name   string
		field  string
		oField string
		value  string
	}{
		{
			name:   "update title",
			field:  "title",
			oField: "Title",
			value:  "New title",
		},
		{
			name:   "update description",
			field:  "description",
			oField: "Description",
			value:  "New description",
		},
		{
			name:   "update time",
			field:  "time",
			oField: "Time",
			value:  time.Now().Add(50 * time.Minute).Format(model.EventDateLayout),
		},
		{
			name:   "update timezone",
			field:  "timezone",
			oField: "Timezone",
			value:  "Europe/Kiev",
		},
	}

	server := mustMakeServer()

	for _, item := range cases {
		t.Run(item.name, func(t *testing.T) {
			event := model.TestEvent(t)
			user := model.TestUser(t)
			rec := httptest.NewRecorder()

			token := getJWTToken(t, user)
			event.UserID = user.ID
			addEventToStore(t, event)

			b := &bytes.Buffer{}
			data := map[string]string{
				item.field: item.value,
			}

			if err := json.NewEncoder(b).Encode(data); err != nil {
				t.Fatalf("could not encode data. %s", err.Error())
			}
			request, _ := http.NewRequest(http.MethodPut, fmt.Sprintf("/api/event/%d", event.ID), b)
			request.Header.Set("Authorization", fmt.Sprintf("Bearer %v", token))

			server.ServeHTTP(rec, request)

			assert.Equal(t, rec.Code, http.StatusCreated)
			assert.Equal(t, rec.Header().Get("content-type"), httpserver.JsonContentType)

			assert.Equal(t, item.value, getFieldValueByName(event, item.oField))
		})
	}
}

func TestServer_HandleDeleteEvent(t *testing.T) {
	t.Run("delete event", func(t *testing.T) {
		server := mustMakeServer()
		event := model.TestEvent(t)
		user := model.TestUser(t)
		rec := httptest.NewRecorder()

		token := getJWTToken(t, user)
		event.UserID = user.ID
		addEventToStore(t, event)

		request, _ := http.NewRequest(http.MethodDelete, fmt.Sprintf("/api/event/%d", event.ID), nil)
		request.Header.Set("Authorization", fmt.Sprintf("Bearer %v", token))

		server.ServeHTTP(rec, request)

		assert.Equal(t, rec.Code, http.StatusOK)
		assert.Equal(t, rec.Header().Get("content-type"), httpserver.JsonContentType)

		_, err := store.Event().FindById(event.ID)
		assert.Error(t, err)
	})
}

func mustMakeServer() *httpserver.Server {
	return httpserver.NewServer(store, jwtWrapper, session)
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
	jb, err := json.Marshal(map[string]bool{
		token: true,
	})

	if err != nil {
		t.Fatal("could not marshal user map session")
	}

	if err := session.Set(strconv.Itoa(user.ID), string(jb)); err != nil {
		t.Fatal("could not save user session")
	}

	return token
}

func addEventToStore(t *testing.T, event *model.Event) {
	if err := store.Event().Create(event); err != nil {
		t.Fatal("could not create event")
	}
}

func getFieldValueByName(o interface{}, field string) string {
	r := reflect.ValueOf(o)
	f := reflect.Indirect(r).FieldByName(field)
	return f.String()
}
