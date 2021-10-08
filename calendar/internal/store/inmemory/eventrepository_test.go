package inmemory_test

import (
	"calendar/internal/model"
	"calendar/internal/store/inmemory"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestEventRepository_Create(t *testing.T) {
	s := inmemory.New()
	e := model.TestEvent(t)

	assert.NoError(t, s.Event().Create(e))
	assert.NotNil(t, e.ID)
}

func TestEventRepository_FindById(t *testing.T) {
	s := inmemory.New()
	e := model.TestEvent(t)
	if err := s.Event().Create(e); err != nil {
		t.Fatal("could not create event")
	}

	r, err := s.Event().FindById(e.ID)

	assert.NoError(t, err)
	assert.NotNil(t, r)
}

func TestEventRepository_Update(t *testing.T) {
	title := "New title"
	s := inmemory.New()
	e := model.TestEvent(t)
	if err := s.Event().Create(e); err != nil {
		t.Fatal("could not create event")
	}

	e.Title = title
	assert.NoError(t, s.Event().Update(e))

	r, err := s.Event().FindById(e.ID)
	assert.NoError(t, err)
	assert.Equal(t, r.Title, title)
}
