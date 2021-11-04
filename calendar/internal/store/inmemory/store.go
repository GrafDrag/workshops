package inmemory

import (
	"calendar/internal/model"
	"calendar/internal/store"
)

type Store struct {
	userRepository  *UserRepository
	eventRepository *EventRepository
}

func New() *Store {
	return &Store{}
}

func (s *Store) User() store.UserRepository {
	if s.userRepository == nil {
		s.userRepository = &UserRepository{
			store: s,
			users: make(map[int]*model.User),
		}
	}

	return s.userRepository
}

func (s *Store) Event() store.EventRepository {
	if s.eventRepository == nil {
		s.eventRepository = &EventRepository{
			store:  s,
			events: make(map[int]*model.Event),
		}
	}

	return s.eventRepository
}
