package store

type Store interface {
	User() UserRepository
	Event() EventRepository
}
