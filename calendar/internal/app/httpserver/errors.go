package httpserver

var (
	errIncorrectLoginOrPassword = "incorrect login or password"
	errEmptyAuthToken           = "missing auth token"
	errInvalidAuthToken         = "invalid/malformed auth token"
	errSessionNotFound          = "session not found"
	errUserNotFound             = "user not found"
	errUserExist                = "this is login use other user"
	errInvalidParams            = "incorrect request params"
	errEventNotFound            = "event not found"
)
