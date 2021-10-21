package store

import "calendar/internal/model"

type UserRepository interface {
	Create(user *model.User) error
	FindById(ID int) (*model.User, error)
	FindByLogin(login string) (*model.User, error)
	Update(*model.User) error
}

type EventRepository interface {
	FindByParams(event model.SearchEvent) ([]*model.Event, error)
	Create(event *model.Event) error
	FindById(ID int) (*model.Event, error)
	Update(*model.Event) error
	Delete(ID int) error
}
