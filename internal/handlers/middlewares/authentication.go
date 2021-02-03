package middlewares

import (
	"encoding/json"
	"net/http"
	"strings"

	"saaa-api/internal/core/config"
	"saaa-api/internal/core/context"
	"saaa-api/internal/core/jwt"
	"saaa-api/internal/models"

	"github.com/go-chi/render"
	"github.com/sirupsen/logrus"
)

const (
	authHeader = "Authorization"
	bearer     = "Bearer "
)

// RequireAuthentication require authentication
func RequireAuthentication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		accessToken := r.Header.Get(authHeader)
		tokenString := strings.Replace(accessToken, bearer, "", 1)
		if tokenString == "" {
			logrus.Error("[RequireAuthentication] token error: ", config.RR.Internal.Unauthorized.Error())
			render.Status(r, config.RR.Internal.Unauthorized.HTTPStatusCode())
			render.JSON(w, r, config.RR.Internal.Unauthorized.WithLocale(r))
			return
		}

		userSession, err := generateUserSessionFromJwtToken(tokenString, true)
		if err != nil {
			logrus.Error("[RequireAuthentication] generate user session from JWT Token: ", err)
			render.Status(r, config.RR.Internal.Unauthorized.HTTPStatusCode())
			render.JSON(w, r, config.RR.Internal.Unauthorized.WithLocale(r))
			return
		}

		context.SetUser(r, userSession)
		next.ServeHTTP(w, r)
	})
}

func generateUserSessionFromJwtToken(token string, onlyValid bool) (*models.UserSession, error) {
	claims, err := jwt.Parsed(token, onlyValid)
	if err != nil {
		return nil, err
	}

	idInterface := claims["sub"]
	var id uint
	if idInterface != nil {
		idByte := []byte(idInterface.(string))
		err = json.Unmarshal(idByte, &id)
		if err != nil {
			return nil, err
		}
	}

	refreshTokenIDInterface := claims["refreshTokenID"]
	var rfTokenID uint
	if refreshTokenIDInterface != nil {
		idByte := []byte(refreshTokenIDInterface.(string))
		err = json.Unmarshal(idByte, &rfTokenID)
		if err != nil {
			return nil, err
		}
	}

	userSession := &models.UserSession{
		ID:             id,
		RefreshTokenID: rfTokenID,
	}

	return userSession, nil
}
