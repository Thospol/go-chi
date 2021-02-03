package healthcheck

import (
	"net/http"

	"saaa-api/internal/core/config"
	"saaa-api/internal/core/result"
	"saaa-api/internal/core/sql"

	"github.com/go-chi/render"
	"github.com/sirupsen/logrus"
)

// Endpoint health check interface
type Endpoint interface {
	HealthCheck() http.HandlerFunc
}

type endpoint struct {
	config *config.Configs
	result *config.ReturnResult
}

// NewEndpoint new endpoint healthCheck
func NewEndpoint(config *config.Configs, result *config.ReturnResult) Endpoint {
	return &endpoint{
		config: config,
		result: result,
	}
}

// HealthCheck endpoint healthCheck
func (ep *endpoint) HealthCheck() http.HandlerFunc {
	handlerFunc := func(w http.ResponseWriter, r *http.Request) {
		sqlDB, err := sql.Database.DB()
		if err != nil {
			logrus.Errorf("[EP-HealthCheck] cannot get database service error: %s", err)
			result.Error(w, r, err)
			return
		}

		err = sqlDB.Ping()
		if err != nil {
			logrus.Errorf("[EP-HealthCheck] call service error: %s", err)
			result.Error(w, r, err)
			return
		}

		render.Status(r, ep.result.Internal.Success.HTTPStatusCode())
		render.JSON(w, r, ep.result.Internal.Success.WithLocale(r))
	}

	return handlerFunc
}
