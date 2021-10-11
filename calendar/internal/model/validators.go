package model

import (
	"errors"
	validation "github.com/go-ozzo/ozzo-validation"
	"time"
)

func RequiredIf(c bool) validation.RuleFunc {
	return func(value interface{}) error {
		if c {
			return validation.Validate(value, validation.Required)
		}

		return nil
	}
}

func TimeZoneValidator(timezone string) validation.RuleFunc {
	return func(value interface{}) error {
		_, err := time.LoadLocation(timezone)
		if err != nil {
			return errors.New("failed validation timezone")
		}
		return nil
	}
}

func DatetimeValidator(layout string, datetime string) validation.RuleFunc {
	return func(value interface{}) error {
		_, err := time.Parse(layout, datetime)
		if err != nil {
			return errors.New("failed validation time")
		}
		return nil
	}
}
