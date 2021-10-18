package sqlstore

import (
	"calendar/internal/model"
	"calendar/internal/store"
	"database/sql"
)

type UserRepository struct {
	store *Store
}

func (r *UserRepository) Create(user *model.User) error {
	if err := user.Validate(); err != nil {
		return err
	}

	if err := user.BeforeCreate(); err != nil {
		return err
	}

	return r.store.db.QueryRow(
		"INSERT INTO users (login, encrypted_password, timezone) VALUES ($1, $2, $3) RETURNING id",
		user.Login,
		user.EncryptedPassword,
		user.Timezone,
	).Scan(&user.ID)
}

func (r *UserRepository) FindById(ID int) (*model.User, error) {
	user := &model.User{}
	err := r.store.db.QueryRow(
		"SELECT id, login, encrypted_password, timezone FROM users WHERE id = $1",
		ID,
	).Scan(&user.ID, &user.Login, &user.EncryptedPassword, &user.Timezone)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.ErrRecordNotFound
		}

		return nil, err
	}

	return user, nil
}

func (r *UserRepository) FindByLogin(login string) (*model.User, error) {
	user := &model.User{}
	err := r.store.db.QueryRow(
		"SELECT id, login, encrypted_password, timezone FROM users WHERE login = $1",
		login,
	).Scan(&user.ID, &user.Login, &user.EncryptedPassword, &user.Timezone)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.ErrRecordNotFound
		}

		return nil, err
	}

	return user, nil
}

func (r *UserRepository) Update(user *model.User) error {
	if err := user.Validate(); err != nil {
		return err
	}

	_, err := r.store.db.Exec(
		"UPDATE users SET login = $1, timezone = $2 WHERE id = $3",
		user.Login,
		user.Timezone,
		user.ID,
	)

	return err
}
