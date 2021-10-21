package sqlstore

import (
	"calendar/internal/model"
	"calendar/internal/store"
	"database/sql"
	"fmt"
	"log"
	"strings"
)

type EventRepository struct {
	store *Store
}

func (r *EventRepository) FindByParams(search model.SearchEvent) ([]*model.Event, error) {
	res := make([]*model.Event, 0)
	var where []string

	where = append(where, fmt.Sprintf("user_id = %d", search.UserID))
	if search.Title != "" {
		where = append(where, fmt.Sprintf("title LIKE '%s'", search.Title))
	}
	if search.Timezone != "" {
		where = append(where, fmt.Sprintf("timezone LIKE '%s'", search.Timezone))
	}

	if search.DateFrom != "" {
		from := search.DateFrom
		if search.TimeFrom != "" {
			from = fmt.Sprintf("%s %s", from, search.TimeFrom)
		} else {
			from = fmt.Sprintf("%s 00:00", from)
		}
		where = append(where, fmt.Sprintf("time > '%s'", from))
	}

	if search.DateTo != "" {
		to := search.DateTo
		if search.TimeTo != "" {
			to = fmt.Sprintf("%s %s", to, search.TimeTo)
		} else {
			to = fmt.Sprintf("%s 00:00", to)
		}
		where = append(where, fmt.Sprintf("time < '%s'", to))
	}

	rows, err := r.store.db.Query(
		"SELECT id, user_id, title, description, time, timezone, duration, notes FROM events WHERE " +
			strings.Join(where, " AND "),
	)

	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			log.Print("failed closing row!")
		}
	}(rows)

	for rows.Next() {
		e := &model.Event{}
		err := rows.Scan(&e.ID, &e.UserID, &e.Title, &e.Description, &e.Time, &e.Timezone, &e.Duration, &e.Notes)
		if err != nil {
			return nil, err
		}
		res = append(res, e)
	}

	return res, nil
}

func (r *EventRepository) Create(event *model.Event) error {
	if err := event.Validate(); err != nil {
		return err
	}

	return r.store.db.QueryRow(
		"INSERT INTO events (user_id, title, description, time, timezone, duration, notes) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id",
		event.UserID,
		event.Title,
		event.Description,
		event.Time,
		event.Timezone,
		event.Duration,
		event.Notes,
	).Scan(&event.ID)
}

func (r *EventRepository) FindById(ID int) (*model.Event, error) {
	e := &model.Event{}

	err := r.store.db.QueryRow(
		"SELECT id, user_id, title, description, time, timezone, duration, notes FROM events WHERE id = $1",
		ID,
	).Scan(&e.ID, &e.UserID, &e.Title, &e.Description, &e.Time, &e.Timezone, &e.Duration, &e.Notes)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.ErrRecordNotFound
		}

		return nil, err
	}

	return e, nil
}

func (r *EventRepository) Update(event *model.Event) error {
	if err := event.Validate(); err != nil {
		return err
	}

	_, err := r.store.db.Exec(
		"UPDATE events SET title = $1, description = $2, time = $3, timezone = $4, duration = $5, notes = $6",
		event.Title,
		event.Description,
		event.Time,
		event.Timezone,
		event.Duration,
		event.Notes,
	)

	return err
}

func (r *EventRepository) Delete(ID int) error {
	_, err := r.store.db.Exec(
		"DELETE FROM events WHERE id = $1",
		ID,
	)

	return err
}
