package controller

import (
	"calendar/internal/auth"
	"calendar/internal/model"
	"calendar/internal/store"
	"context"
	"net/http"
	"strconv"
)

type EventController struct {
	store store.Store
}

func NewEventController(store store.Store) *EventController {
	return &EventController{
		store: store,
	}
}

func (c *EventController) List(ctx context.Context, searchModel model.SearchEvent) ([]*model.Event, *ResponseError) {
	user, err := c.store.User().FindById(ctx.Value(auth.KeyUserID).(int))
	if err != nil {
		return nil, &ResponseError{
			status: http.StatusForbidden,
			Err:    errUserNotFound,
		}
	}

	searchModel.UserID = user.ID

	res, err := c.store.Event().FindByParams(searchModel)

	if err != nil {
		return nil, &ResponseError{
			Err: err,
		}
	}

	return res, nil
}

func (c *EventController) FindById(ctx context.Context, objId interface{}) (*model.Event, *ResponseError) {
	id, err := getID(objId)
	if err != nil {
		return nil, &ResponseError{
			status: http.StatusBadRequest,
			Err:    err,
		}
	}

	event, err := c.store.Event().FindById(id)
	if err != nil || event.UserID != ctx.Value(auth.KeyUserID).(int) {
		return nil, &ResponseError{
			status: http.StatusNotFound,
			Err:    errEventNotFound,
		}
	}

	return event, nil
}

func (c *EventController) Create(ctx context.Context, event *model.Event) *ResponseError {
	user, err := c.store.User().FindById(ctx.Value(auth.KeyUserID).(int))
	if err != nil {
		return &ResponseError{
			status: http.StatusForbidden,
			Err:    errUserNotFound,
		}
	}

	event.UserID = user.ID

	if err := event.Validate(); err != nil {
		return &ResponseError{
			status: http.StatusBadRequest,
			Err:    err,
		}
	}

	if err := c.store.Event().Create(event); err != nil {
		return &ResponseError{
			status: http.StatusInternalServerError,
			Err:    err,
		}
	}

	return nil
}

func (c *EventController) Update(ctx context.Context, event *model.Event) *ResponseError {
	repository := c.store.Event()
	res, err := repository.FindById(event.ID)
	if err != nil || res.UserID != ctx.Value(auth.KeyUserID).(int) {
		return &ResponseError{
			status: http.StatusNotFound,
			Err:    errEventNotFound,
		}
	}

	event.UserID = res.UserID

	if err := event.Validate(); err != nil {
		return &ResponseError{
			status: http.StatusBadRequest,
			Err:    err,
		}
	}

	if err = repository.Update(event); err != nil {
		return &ResponseError{
			status: http.StatusInternalServerError,
			Err:    err,
		}
	}

	return nil
}

func (c *EventController) Delete(ctx context.Context, objId interface{}) *ResponseError {
	repository := c.store.Event()
	id, err := getID(objId)
	if err != nil {
		return &ResponseError{
			status: http.StatusBadRequest,
			Err:    err,
		}
	}

	event, err := repository.FindById(id)
	if err != nil || event.UserID != ctx.Value(auth.KeyUserID).(int) {
		return &ResponseError{
			status: http.StatusNotFound,
			Err:    errEventNotFound,
		}
	}

	if err := repository.Delete(event.ID); err != nil {
		return &ResponseError{
			status: http.StatusNotFound,
			Err:    errEventNotFound,
		}
	}

	return nil
}

func getID(objId interface{}) (int, error) {
	if id, ok := objId.(int); ok {
		return id, nil
	}

	if id, ok := objId.(int32); ok {
		return int(id), nil
	}

	if id, ok := objId.(string); ok {
		id, err := strconv.Atoi(id)
		if err != nil {
			return 0, err
		}
		return id, nil
	}

	return 0, errInvalidParams
}