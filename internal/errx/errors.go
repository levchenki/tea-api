package errx

import (
	"net/http"
)

type ErrorResponse struct {
	Err            error  `json:"-"`
	HTTPStatusCode int    `json:"-"`
	ErrorText      string `json:"error,omitempty"`
}

func (e *ErrorResponse) Error() string {
	return e.Err.Error()
}

func ErrorBadRequest(err error) *ErrorResponse {
	e := &ErrorResponse{
		Err:            err,
		HTTPStatusCode: http.StatusBadRequest,
		ErrorText:      err.Error(),
	}
	return e
}

func ErrorInternalServer(err error) *ErrorResponse {
	e := &ErrorResponse{
		Err:            err,
		HTTPStatusCode: http.StatusInternalServerError,
		ErrorText:      err.Error(),
	}
	return e
}

func ErrorNotFound(err error) *ErrorResponse {
	e := &ErrorResponse{
		Err:            err,
		HTTPStatusCode: http.StatusNotFound,
		ErrorText:      err.Error(),
	}
	return e
}

func ErrorUnauthorized(err error) *ErrorResponse {
	e := &ErrorResponse{
		Err:            err,
		HTTPStatusCode: http.StatusUnauthorized,
		ErrorText:      err.Error(),
	}
	return e
}

func ErrorForbidden(err error) *ErrorResponse {
	e := &ErrorResponse{
		Err:            err,
		HTTPStatusCode: http.StatusForbidden,
		ErrorText:      err.Error(),
	}
	return e
}
