package result

import (
	"net/http"

	"saaa-api/internal/core/config"
	"saaa-api/internal/core/context"

	"github.com/go-chi/render"
)

// Error render error to client
func Error(w http.ResponseWriter, r *http.Request, err error) {
	errMsg := config.RR.Internal.ConnectionError
	if locErr, ok := err.(config.Result); ok {
		errMsg = locErr
	}
	context.SetErrMsg(r, errMsg.Error())
	render.Status(r, errMsg.HTTPStatusCode())
	render.JSON(w, r, errMsg.WithLocale(r))
}
