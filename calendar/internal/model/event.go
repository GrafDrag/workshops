package model

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"github.com/go-ozzo/ozzo-validation/v4"
)

const EventDateLayout = "2006-01-02 15:04"

type Event struct {
	ID          int    `json:"id"`
	UserID      int    `json:"-"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Time        string `json:"time"`
	Timezone    string `json:"timezone"`
	Duration    int32  `json:"duration"`
	Notes       Node   `json:"notes,omitempty"`
}

type Node []string

func (n Node) Value() (driver.Value, error) {
	return json.Marshal(n)
}

func (n *Node) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(b, &n)
}

func (e Event) Validate() error {
	return validation.ValidateStruct(
		&e,
		validation.Field(&e.UserID, validation.Required),
		validation.Field(&e.Title, validation.Required, validation.Length(1, 150)),
		validation.Field(&e.Description, validation.Required, validation.Length(5, 1000)),

		validation.Field(&e.Time, validation.By(DatetimeValidator(EventDateLayout, e.Time))),
		validation.Field(&e.Timezone, validation.By(TimeZoneValidator(e.Timezone))),

		validation.Field(&e.Notes, validation.Each(validation.Length(5, 50))),
	)
}
