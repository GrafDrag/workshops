package model_test

import (
	"calendar/internal/model"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestEvent_Validate(t *testing.T) {
	cases := []struct {
		name    string
		e       func() *model.Event
		isValid bool
	}{
		{
			name: "valid event",
			e: func() *model.Event {
				return model.TestEvent(t)
			},
			isValid: true,
		},
		{
			name: "empty title",
			e: func() *model.Event {
				e := model.TestEvent(t)
				e.Title = ""
				return e
			},
			isValid: false,
		},
		{
			name: "empty description",
			e: func() *model.Event {
				e := model.TestEvent(t)
				e.Description = ""
				return e
			},
			isValid: false,
		},
		{
			name: "empty time",
			e: func() *model.Event {
				e := model.TestEvent(t)
				e.Time = ""
				return e
			},
			isValid: false,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.isValid {
				assert.NoError(t, tc.e().Validate())
			} else {
				assert.Error(t, tc.e().Validate())
			}
		})
	}
}
