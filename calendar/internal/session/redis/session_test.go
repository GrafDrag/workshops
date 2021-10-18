package redis_test

import (
	"calendar/internal/session/redis"
	"github.com/stretchr/testify/assert"
	"testing"
)

var (
	sessionKey   = "test_key"
	sessionValue = "test_value"
	configPath   = "../../../configs/httpserver.toml"
)

func TestSession_Get(t *testing.T) {
	s := redis.TestSession(t, configPath)
	if err := s.Set(sessionKey, sessionValue); err != nil {
		t.Fatalf("faled set to session, %v", err)
	}

	v, err := s.Get(sessionKey)

	assert.NoError(t, err)
	assert.Equal(t, sessionValue, v)
}

func TestSession_Set(t *testing.T) {
	s := redis.TestSession(t, configPath)
	err := s.Set(sessionKey, sessionValue)

	assert.NoError(t, err)
}

func TestSession_Remove(t *testing.T) {
	s := redis.TestSession(t, configPath)
	if err := s.Set(sessionKey, sessionValue); err != nil {
		t.Fatalf("faled set to session, %v", err)
	}

	v, err := s.Remove(sessionKey)

	assert.NoError(t, err)
	assert.Equal(t, sessionValue, v)

	_, err = s.Get(sessionKey)

	assert.Error(t, err)
}

func TestSession_Flash(t *testing.T) {
	s := redis.TestSession(t, configPath)
	if err := s.Set(sessionKey, sessionValue); err != nil {
		t.Fatalf("faled set to session, %v", err)
	}

	assert.NoError(t, s.Flash())

	_, err := s.Get(sessionKey)

	assert.Error(t, err)
}
