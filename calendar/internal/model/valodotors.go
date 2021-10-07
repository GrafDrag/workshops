package model

import (
	"errors"
	validation "github.com/go-ozzo/ozzo-validation"
	"time"
)

func requiredIf(c bool) validation.RuleFunc {
	return func(value interface{}) error {
		if c {
			return validation.Validate(value, validation.Required)
		}

		return nil
	}
}

func timeZoneValidator(timezone string) validation.RuleFunc {
	return func(value interface{}) error {
		_, err := time.LoadLocation(timezone)
		if err != nil {
			return errors.New("failed validation timezone")
		}
		return nil
	}
}

func datetimeValidator(datetime string) validation.RuleFunc {
	return func(value interface{}) error {
		_, err := time.Parse("2006-01-02 15:04:05", datetime)
		if err != nil {
			return errors.New("failed validation time")
		}
		return nil
	}
}
