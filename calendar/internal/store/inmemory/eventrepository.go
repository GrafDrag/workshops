package inmemory

import (
	"calendar/internal/model"
	"calendar/internal/store"
)

type EventRepository struct {
	store  *Store
	events map[int]*model.Event
}

func (r EventRepository) FindByUser(userID int) []*model.Event {
	var res []*model.Event

	for _, event := range r.events {
		if event.UserID == userID {
			res = append(res, event)
		}
	}

	return res
}

func (r EventRepository) Create(event *model.Event) error {
	err := event.Validate()
	if err != nil {
		return err
	}

	event.ID = len(r.events) + 1
	r.events[event.ID] = event

	return nil
}

func (r EventRepository) FindById(ID int) (*model.Event, error) {
	event, ok := r.events[ID]
	if !ok {
		return nil, store.ErrRecordNotFound
	}
	return event, nil
}

func (r EventRepository) Update(event *model.Event) error {
	err := event.Validate()
	if err != nil {
		return err
	}

	if _, ok := r.events[event.ID]; !ok {
		return store.ErrRecordNotFound
	}

	r.events[event.ID] = event

	return nil
}
