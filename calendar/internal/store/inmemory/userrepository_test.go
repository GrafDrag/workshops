package inmemory_test

import (
	"calendar/internal/model"
	"calendar/internal/store/inmemory"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUserRepository_Create(t *testing.T) {
	s := inmemory.New()
	u := model.TestUser(t)

	assert.NoError(t, s.User().Create(u))
	assert.NotNil(t, u.ID)
}

func TestUserRepository_FindById(t *testing.T) {
	s := inmemory.New()
	u := model.TestUser(t)
	if err := s.User().Create(u); err != nil {
		t.Fatal("could not create user")
	}

	r, err := s.User().FindById(u.ID)

	assert.NoError(t, err)
	assert.NotNil(t, r)
}

func TestUserRepository_FindByLogin(t *testing.T) {
	s := inmemory.New()
	u := model.TestUser(t)
	if err := s.User().Create(u); err != nil {
		t.Fatal("could not create user")
	}

	r, err := s.User().FindByLogin(u.Login)

	assert.NoError(t, err)
	assert.NotNil(t, r)
}

func TestUserRepository_Update(t *testing.T) {
	tz := "Europe/Kiev"
	s := inmemory.New()
	u := model.TestUser(t)
	if err := s.User().Create(u); err != nil {
		t.Fatal("could not create user")
	}

	u.Timezone = tz
	assert.NoError(t, s.User().Update(u))

	r, err := s.User().FindByLogin(u.Login)
	assert.NoError(t, err)
	assert.Equal(t, r.Timezone, tz)
}
