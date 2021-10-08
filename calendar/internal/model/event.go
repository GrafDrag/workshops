package model

import (
	validation "github.com/go-ozzo/ozzo-validation"
)

type Event struct {
	ID          int
	UserID      int
	Title       string `json:"title"`
	Description string `json:"description"`
	Time        string `json:"time"`
	Timezone    string `json:"timezone"`
	Duration    int32  `json:"duration"`
	Notes       []string
}

func (e Event) Validate() error {
	return validation.ValidateStruct(
		&e,
		validation.Field(&e.UserID, validation.Required),
		validation.Field(&e.Title, validation.Required, validation.Length(1, 150)),
		validation.Field(&e.Description, validation.Required, validation.Length(5, 1000)),

		validation.Field(&e.Time, validation.By(DatetimeValidator(e.Time))),
		validation.Field(&e.Timezone, validation.By(TimeZoneValidator(e.Timezone))),

		validation.Field(&e.Notes, validation.Each(validation.Length(5, 50))),
	)
}
