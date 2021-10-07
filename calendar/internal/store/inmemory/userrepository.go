package inmemory

import (
	"calendar/internal/model"
	"calendar/internal/store"
)

type UserRepository struct {
	store *Store
	users map[string]*model.User
}

func (r UserRepository) Create(user *model.User) error {
	err := user.Validate()
	if err != nil {
		return err
	}

	err = user.BeforeCreate()
	if err != nil {
		return err
	}

	r.users[user.Login] = user
	user.ID = len(r.users)

	return nil
}

func (r UserRepository) FindById(ID int) (*model.User, error) {
	if ID > len(r.users) {
		return nil, store.ErrRecordNotFound
	}
	i := 0
	for _, user := range r.users {
		i++
		if i == ID {
			return user, nil
		}
	}
	return nil, store.ErrRecordNotFound
}

func (r UserRepository) FindByLogin(login string) (*model.User, error) {
	user, ok := r.users[login]
	if !ok {
		return nil, store.ErrRecordNotFound
	}
	return user, nil
}

func (r UserRepository) Update(user *model.User) error {
	err := user.Validate()
	if err != nil {
		return err
	}

	if _, ok := r.users[user.Login]; !ok {
		return store.ErrRecordNotFound
	}

	r.users[user.Login] = user

	return nil
}
