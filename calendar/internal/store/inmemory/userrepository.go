package inmemory

import (
	"calendar/internal/model"
	"calendar/internal/store"
)

type UserRepository struct {
	store *Store
	users map[int]*model.User
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
	user.ID = len(r.users) + 1
	r.users[user.ID] = user

	return nil
}

func (r UserRepository) FindById(ID int) (*model.User, error) {
	user, ok := r.users[ID]
	if !ok {
		return nil, store.ErrRecordNotFound
	}

	return user, nil
}

func (r UserRepository) FindByLogin(login string) (*model.User, error) {
	for _, user := range r.users {
		if user.Login == login {
			return user, nil
		}
	}

	return nil, store.ErrRecordNotFound
}

func (r UserRepository) Update(user *model.User) error {
	err := user.Validate()
	if err != nil {
		return err
	}

	if _, ok := r.users[user.ID]; !ok {
		return store.ErrRecordNotFound
	}

	r.users[user.ID] = user

	return nil
}
