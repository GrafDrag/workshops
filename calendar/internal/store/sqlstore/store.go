package sqlstore

import (
	"calendar/internal/store"
	"database/sql"
)

type Store struct {
	db *sql.DB

	userRepository  *UserRepository
	eventRepository *EventRepository
}

func New(db *sql.DB) *Store {
	return &Store{
		db: db,
	}
}

func (s *Store) User() store.UserRepository {
	if s.userRepository == nil {
		s.userRepository = &UserRepository{
			store: s,
		}
	}

	return s.userRepository
}

func (s *Store) Event() store.EventRepository {
	if s.eventRepository == nil {
		s.eventRepository = &EventRepository{
			store: s,
		}
	}

	return s.eventRepository
}
