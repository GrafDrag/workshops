package sqlstore

import (
	"database/sql"
	"fmt"
	"github.com/BurntSushi/toml"
	"strings"
	"testing"

	_ "github.com/lib/pq"
)

type conf struct {
	DB database `toml:"database"`
}

type database struct {
	Host     string
	Port     int
	User     string
	Password string
	DB       string
}

func TestDB(t *testing.T, configPath string) (*sql.DB, func(...string)) {
	t.Helper()

	config := &conf{}
	_, err := toml.DecodeFile(configPath, config)
	if err != nil {
		t.Fatal("failed read config file", err)
	}

	psqlconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s_test sslmode=disable", config.DB.Host, config.DB.Port, config.DB.User, config.DB.Password, config.DB.DB)
	db, err := sql.Open("postgres", psqlconn)
	if err != nil {
		t.Fatal("failed create connection", err)
	}

	if err := db.Ping(); err != nil {
		t.Fatal("failed ping to database", err)
	}

	return db, func(tables ...string) {
		if len(tables) > 0 {
			q := fmt.Sprintf("TRUNCATE %s CASCADE", strings.Join(tables, ", "))
			if _, err := db.Exec(q); err != nil {
				t.Fatal("failed not clear database", err)
			}
		}

		db.Close()
	}
}
