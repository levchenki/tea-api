package errx

import "net/http"

type AppError struct {
	Code    int    `json:"-"`
	Message string `json:"error,omitempty"`
}

func (e *AppError) Error() string {
	return e.Message
}

func NewAppError(code int, error error) *AppError {
	return &AppError{
		Code:    code,
		Message: error.Error(),
	}
}

func NewBadRequestError(error error) *AppError {
	return &AppError{
		Code:    http.StatusBadRequest,
		Message: error.Error(),
	}
}
func NewNotFoundError(error error) *AppError {
	return &AppError{
		Code:    http.StatusNotFound,
		Message: error.Error(),
	}
}

func NewForbiddenError(error error) *AppError {
	return &AppError{
		Code:    http.StatusForbidden,
		Message: error.Error(),
	}
}

func NewUnauthorizedError(error error) *AppError {
	return &AppError{
		Code:    http.StatusUnauthorized,
		Message: error.Error(),
	}
}

func NewInternalServerError(error error) *AppError {
	return &AppError{
		Code:    http.StatusInternalServerError,
		Message: error.Error(),
	}
}
