package sqlstore

import (
	"calendar/internal/config"
	"calendar/internal/store"
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
)

type Store struct {
	db *sql.DB

	userRepository  *UserRepository
	eventRepository *EventRepository
}

func New(config config.DBConfig) (*Store, func() error, error) {
	db, err := newDB(config.Host, config.Port, config.User, config.Password, config.DB)
	if err != nil {
		return nil, nil, err
	}

	return &Store{
		db: db,
	}, db.Close, nil
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

func newDB(host string, port int, user string, password string, dbname string) (*sql.DB, error) {
	conn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	db, err := sql.Open("postgres", conn)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
