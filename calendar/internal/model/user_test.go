package model_test

import (
	"calendar/internal/model"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUser_Validate(t *testing.T) {
	cases := []struct {
		name    string
		u       func() *model.User
		isValid bool
	}{
		{
			name: "valid user",
			u: func() *model.User {
				return model.TestUser(t)
			},
			isValid: true,
		},
		{
			name: "valid user with encrypt password",
			u: func() *model.User {
				u := model.TestUser(t)
				u.Password = ""
				u.EncryptedPassword = "not_empty"
				return u
			},
			isValid: true,
		},
		{
			name: "empty login",
			u: func() *model.User {
				u := model.TestUser(t)
				u.Login = ""
				return u
			},
			isValid: false,
		},
		{
			name: "invalid Timezone",
			u: func() *model.User {
				u := model.TestUser(t)
				u.Timezone = "test"
				return u
			},
			isValid: false,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.isValid {
				assert.NoError(t, tc.u().Validate())
			} else {
				assert.Error(t, tc.u().Validate())
			}
		})
	}
}

func TestUser_BeforeCreate(t *testing.T) {
	u := model.TestUser(t)

	assert.NoError(t, u.BeforeCreate())
	assert.NotNil(t, u.EncryptedPassword)
}

func TestUser_CheckPassword(t *testing.T) {
	p := "password"
	u := model.TestUser(t)

	u.Password = p
	u.BeforeCreate()
	assert.NoError(t, u.CheckPassword(p))
}
