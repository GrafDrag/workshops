package inmemory

import (
	"calendar/internal/model"
	"calendar/internal/store"
	"fmt"
	"strings"
	"time"
)

type EventRepository struct {
	store  *Store
	events map[int]*model.Event
}

func (r EventRepository) FindByParams(search model.SearchEvent) ([]*model.Event, error) {
	res := make([]*model.Event, 0)

	for _, event := range r.events {
		if event.UserID != search.UserID {
			continue
		}
		if search.Title != "" && !strings.Contains(event.Title, search.Title) {
			continue
		}
		if search.Timezone != "" && event.Timezone != search.Timezone {
			continue
		}
		eventTime, err := time.Parse(model.EventDateLayout, event.Time)

		if err != nil {
			continue
		}

		if search.DateFrom != "" {
			from := search.DateFrom
			if search.TimeFrom != "" {
				from = fmt.Sprintf("%s %s", from, search.TimeFrom)
			} else {
				from = fmt.Sprintf("%s 00:00", from)
			}
			fromDate, err := time.Parse(model.EventDateLayout, from)
			if err != nil || eventTime.Before(fromDate) {
				continue
			}
		}

		if search.DateTo != "" {
			to := search.DateTo
			if search.TimeTo != "" {
				to = fmt.Sprintf("%s %s", to, search.TimeTo)
			} else {
				to = fmt.Sprintf("%s 00:00", to)
			}
			toDate, err := time.Parse(model.EventDateLayout, to)
			if err != nil || eventTime.After(toDate) {
				continue
			}
		}

		res = append(res, event)
	}

	return res, nil
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

func (r EventRepository) Delete(ID int) error {
	if _, ok := r.events[ID]; !ok {
		return store.ErrRecordNotFound
	}

	delete(r.events, ID)

	return nil
}
