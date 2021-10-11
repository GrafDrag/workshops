package model

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID                int
	Login             string `json:"login"`
	Password          string `json:"password"`
	EncryptedPassword string
	Timezone          string `json:"timezone"`
}

func (u User) Validate() error {
	return validation.ValidateStruct(
		&u,
		validation.Field(&u.Login, validation.Required, validation.Length(6, 20).Error("User login invalid")),
		validation.Field(&u.Password, validation.By(RequiredIf(u.EncryptedPassword == "")), validation.Length(6, 20)),
		validation.Field(&u.Timezone, validation.By(TimeZoneValidator(u.Timezone))),
	)
}

func (u *User) BeforeCreate() error {
	if len(u.Password) > 0 {
		enc, err := encryptString(u.Password)
		if err != nil {
			return err
		}

		u.EncryptedPassword = enc
	}

	return nil
}

func (u *User) CheckPassword(password string) error {
	if err := bcrypt.CompareHashAndPassword([]byte(u.EncryptedPassword), []byte(password)); err != nil {
		return err
	}

	return nil
}

func encryptString(s string) (string, error) {
	b, err := bcrypt.GenerateFromPassword([]byte(s), bcrypt.MinCost)
	if err != nil {
		return "", err
	}

	return string(b), nil
}
