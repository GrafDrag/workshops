package httpserver

import (
	"calendar/internal/app/httpserver/auth"
	redis2 "calendar/internal/session/redis"
	"calendar/internal/store/sqlstore"
	"context"
	"database/sql"
	"fmt"
	"github.com/go-redis/redis/v8"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"net/http"
)

var server *Server

func Start(config *Config) error {
	db, err := newDB(config.DB.Host, config.DB.Port, config.DB.User, config.DB.Password, config.DB.DB)
	if err != nil {
		return err
	}

	defer db.Close()

	store := sqlstore.New(db)
	wrapper := &auth.JwtWrapper{
		SecretKey:       config.Jwt.JwtSecretKey,
		ExpirationHours: config.Jwt.JwtExpHours,
		Issuer:          "AuthService",
	}

	rdb, err := newSession(config.Session.Host, config.Session.Port, config.Session.Password, config.Session.DB)
	if err != nil {
		return err
	}
	session := redis2.NewSession(rdb)
	server = NewServer(store, wrapper, session)

	logrus.Infof("Server starting on %v...", config.BindAddr)
	return http.ListenAndServe(config.BindAddr, server)
}

func newDB(host string, port int, user string, password string, dbname string) (*sql.DB, error) {
	psqlconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	db, err := sql.Open("postgres", psqlconn)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}

func newSession(host string, port int, password string, db int) (*redis.Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", host, port),
		Password: password, // no password set
		DB:       db,       // use default DB
	})

	err := rdb.Ping(context.Background()).Err()
	if err != nil {
		return nil, err
	}

	return rdb, nil
}
