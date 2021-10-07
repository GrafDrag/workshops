package store

import "calendar/internal/model"

type UserRepository interface {
	Create(user *model.User) error
	FindById(ID int) (*model.User, error)
	FindByLogin(login string) (*model.User, error)
	Update(*model.User) error
}

type EventRepository interface {
	FindByUser(userID int) []*model.Event
	Create(user *model.Event) error
	FindById(ID int) (*model.Event, error)
	Update(*model.Event) error
}
