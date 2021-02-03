package guest

import (
	"net/http"

	"saaa-api/internal/core/config"
	"saaa-api/internal/core/result"
	"saaa-api/internal/core/sql"
	"saaa-api/internal/core/validate"
	"saaa-api/internal/pkg/token"

	"github.com/go-chi/render"
	"github.com/sirupsen/logrus"
)

// Endpoint endpoint guest interface
type Endpoint interface {
	RefreshToken() http.HandlerFunc
}

type endpoint struct {
	config       *config.Configs
	result       *config.ReturnResult
	tokenService token.Service
}

// NewEndpoint new endpoint guest
func NewEndpoint(config *config.Configs, result *config.ReturnResult) Endpoint {
	return &endpoint{
		config:       config,
		result:       result,
		tokenService: token.NewService(config, result),
	}
}

// RefreshToken godoc
// @Tags Authentication
// @Summary RefreshToken
// @Description Request RefreshToken
// @Accept json
// @Produce json
// @Param Accept-Language header string false "(en, th)"
// @Param request body guest.refreshTokenRequest true "request body"
// @Success 200 {object} models.Token
// @Failure 400 {object} config.SwaggerInfoResult
// @Security ApiKeyAuth
// @Router /c/guest/refreshToken [post]
func (ep *endpoint) RefreshToken() http.HandlerFunc {
	handlerFunc := func(w http.ResponseWriter, r *http.Request) {
		var request refreshTokenRequest
		validate.BindData(w, r, &request, func() {
			validate.GetDatabase(w, r, func(mysqlConf *sql.MysqlConfig) {
				response, err := ep.tokenService.RefreshToken(mysqlConf.Database, request.RefreshToken)
				if err != nil {
					logrus.Errorf("[EP-RefreshToken] call service refresh token error: %s", err)
					result.Error(w, r, err)
					return
				}

				render.Status(r, config.RR.Internal.Success.HTTPStatusCode())
				render.JSON(w, r, response)
			})
		})
	}

	return handlerFunc
}
