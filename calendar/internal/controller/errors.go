package controller

import (
	"errors"
	"net/http"
)

var (
	errIncorrectLoginOrPassword = errors.New("incorrect login or password")
	errUserNotFound             = errors.New("user not found")
	errUserExist                = errors.New("this is login use other user")
	errInvalidParams            = errors.New("incorrect request params")
	errEventNotFound            = errors.New("event not found")
)

type ResponseError struct {
	status int
	Err    error
}

func (e *ResponseError) Error() string {
	return e.Err.Error()
}

func (e *ResponseError) GetStatus() int {
	if e.status == 0 {
		e.status = http.StatusInternalServerError
	}
	return e.status
}
