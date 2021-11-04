package model

import (
	"testing"
	"time"
)

func TestUser(t *testing.T) *User {
	t.Helper()

	return &User{
		Login:    "testing",
		Password: "password",
		Timezone: time.Now().Location().String(),
	}
}

func TestEvent(t *testing.T) *Event {
	t.Helper()

	return &Event{
		UserID:      1,
		Title:       "test",
		Description: "Test description",
		Time:        time.Now().Add(10 * time.Minute).Format("2006-01-02 15:04"),
		Timezone:    time.Now().Location().String(),
		Duration:    5,
		Notes: []string{
			"test1", "test2",
		},
	}
}
