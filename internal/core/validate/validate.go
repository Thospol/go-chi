package validate

import (
	"errors"
	"net/http"

	"saaa-api/internal/core/bind"
	"saaa-api/internal/core/config"
	"saaa-api/internal/core/context"
	"saaa-api/internal/core/sql"
	"saaa-api/internal/models"

	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
)

// UserSession validate user session
func UserSession(w http.ResponseWriter, r *http.Request, callBack func(*models.UserSession)) {
	userSession, ok := context.GetUser(r)
	if !ok {
		logrus.Error("[UserSession] error: usersession expired")
		render.Status(r, config.RR.Internal.Unauthorized.HTTPStatusCode())
		render.JSON(w, r, config.RR.Internal.Unauthorized.WithLocale(r))
		return
	}

	callBack(userSession)
}

// BindDataUserSession validate data and user session
func BindDataUserSession(w http.ResponseWriter, r *http.Request, request interface{}, callBack func(*models.UserSession)) {
	userSession, ok := context.GetUser(r)
	if !ok {
		logrus.Error("[BindDataUserSession] validate data and user session error: usersession expired")
		render.Status(r, config.RR.Internal.Unauthorized.HTTPStatusCode())
		render.JSON(w, r, config.RR.Internal.Unauthorized.WithLocale(r))
		return
	}

	bind.Bind(r, request)
	context.SetParameters(r, request)

	if err := config.CF.Validator.Struct(request); err != nil {
		logrus.Errorf("[BindDataUserSession] validate request error: %s", err)
		errs := err.(validator.ValidationErrors)
		lang := config.RR.GetLanguage(r)
		trans, _ := config.CF.UniversalTranslator.GetTranslator(lang)

		for _, e := range errs {
			err = errors.New(e.Translate(trans))
		}

		render.Status(r, config.RR.Internal.BadRequest.HTTPStatusCode())
		render.JSON(w, r, config.RR.CustomMessage(err.Error(), err.Error()).WithLocale(r))
		return
	}

	callBack(userSession)
}

// BindData validate data
func BindData(w http.ResponseWriter, r *http.Request, request interface{}, callBack func()) {
	bind.Bind(r, request)
	context.SetParameters(r, request)

	if err := config.CF.Validator.Struct(request); err != nil {
		logrus.Errorf("[BindData] validate request error: %s", err)
		render.Status(r, config.RR.Internal.BadRequest.HTTPStatusCode())
		render.JSON(w, r, config.RR.CustomMessage(err.Error(), err.Error()).WithLocale(r))
		return
	}

	callBack()
}

// GetDatabase get database
func GetDatabase(w http.ResponseWriter, r *http.Request, callBack func(*sql.MysqlConfig)) {
	sqlConf, ok := context.GetDatabase(r)
	if !ok {
		logrus.Errorf("[GetDatabase] error: get database from context")
		render.Status(r, config.RR.Internal.DatabaseNotFound.HTTPStatusCode())
		render.JSON(w, r, config.RR.Internal.DatabaseNotFound.WithLocale(r))
		return
	}

	callBack(sqlConf)
}
