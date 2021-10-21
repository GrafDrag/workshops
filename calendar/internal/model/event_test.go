package model_test

import (
	"calendar/internal/model"
	"encoding/json"
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

func TestNode_Value(t *testing.T) {
	t.Run("check get value", func(t *testing.T) {
		n := model.Node{}

		value, err := n.Value()
		marshal, _ := json.Marshal(n)

		assert.NoError(t, err)
		assert.Equal(t, value, marshal)
	})
}

func TestNode_Scan(t *testing.T) {
	t.Run("check get value", func(t *testing.T) {
		n := model.Node{"test"}
		marshal, _ := json.Marshal(n)

		var s model.Node
		err := s.Scan(marshal)

		assert.NoError(t, err)
		assert.Equal(t, s, n)
	})
}
