package redis

import (
	"context"
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/go-redis/redis/v8"
	"testing"

	_ "github.com/lib/pq"
)

type conf struct {
	Session sessionConfig `toml:"session"`
}

type sessionConfig struct {
	Host     string
	Port     int
	Password string
	DB       int
}

func TestSession(t *testing.T, configPath string) *session {
	t.Helper()

	config := &conf{}
	_, err := toml.DecodeFile(configPath, config)
	if err != nil {
		t.Fatal("failed read config file", err)
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", config.Session.Host, config.Session.Port),
		Password: config.Session.Password, // no password set
		DB:       1,                       // use default DB
	})

	err = rdb.Ping(context.Background()).Err()
	if err != nil {
		t.Fatal("failed connect to redis", err)
	}

	return &session{
		redis: rdb,
	}
}
