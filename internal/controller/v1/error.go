package v1

import (
	"errors"
	"github.com/go-chi/render"
	"github.com/levchenki/tea-api/internal/errx"
	"net/http"
)

func handleError(w http.ResponseWriter, r *http.Request, err error) {
	var appError *errx.AppError
	if !errors.As(err, &appError) {
		appError = errx.NewAppError(http.StatusInternalServerError, err)
	}
	render.Status(r, appError.Code)
	render.JSON(w, r, appError)
}
