package session

type Session interface {
	Get(key string) (string, error)
	Set(key string, value string) error
	Remove(key string) (string, error)
	Flash() error
}
