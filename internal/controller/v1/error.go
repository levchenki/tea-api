package v1

import (
	"errors"
	"fmt"
	"github.com/go-chi/render"
	"github.com/levchenki/tea-api/internal/errx"
	"github.com/levchenki/tea-api/internal/logx"
	"net/http"
)

func handleError(w http.ResponseWriter, r *http.Request, log logx.AppLogger, err error) {
	var appError *errx.AppError
	if !errors.As(err, &appError) {
		appError = errx.NewAppError(http.StatusInternalServerError, err)
	}

	logMessage := fmt.Sprintf("Error occurred: %s | path=%s | method=%s", appError.Error(), r.URL.Path, r.Method)
	log.Error(logMessage)

	render.Status(r, appError.Code)
	render.JSON(w, r, appError)
}
