package token

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"saaa-api/internal/core/config"
	"saaa-api/internal/core/jwt"
	"saaa-api/internal/models"
	"saaa-api/internal/repositories"
	"strconv"
	"time"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// Service service token interface
type Service interface {
	GenerateAccessToken(database *gorm.DB, userID uint) (*models.Token, error)
	RefreshToken(database *gorm.DB, refreshTokenString string) (*models.Token, error)
}

type service struct {
	config          *config.Configs
	result          *config.ReturnResult
	tokenRepository repositories.RefreshTokenRepository
}

// NewService new service token
func NewService(config *config.Configs, result *config.ReturnResult) Service {
	return &service{
		config:          config,
		result:          result,
		tokenRepository: repositories.NewRefreshTokenRepository(),
	}
}

func (s *service) generateRefreshToken(userID string) string {
	hasher := sha256.New()
	_, _ = hasher.Write([]byte(fmt.Sprintf("%s_%s", userID, time.Now().String())))
	sha := base64.URLEncoding.EncodeToString(hasher.Sum(nil))
	return sha
}

// GenerateAccessToken generate access token with jwt
func (s *service) GenerateAccessToken(database *gorm.DB, userID uint) (*models.Token, error) {
	tn := time.Now()
	expTime := tn.Add(s.getExpireTime())
	refreshTokenString := s.generateRefreshToken(strconv.FormatUint(uint64(userID), 10))
	refreshToken := &models.RefreshToken{
		RefreshToken: refreshTokenString,
		ExpireTime:   expTime,
		Used:         false,
		UserID:       userID,
	}
	err := s.tokenRepository.Create(database, refreshToken)
	if err != nil {
		logrus.Errorf("[GenerateAccessToken] create refreshToken error: %s", err)
		return nil, err
	}

	idByte, err := json.Marshal(userID)
	if err != nil {
		logrus.Errorf("[GenerateAccessToken] marshal userID error: %s", err)
		return nil, err
	}

	refreshTokenIDByte, err := json.Marshal(refreshToken.ID)
	if err != nil {
		logrus.Errorf("[GenerateAccessToken] marshal refreshTokenID error: %s", err)
		return nil, err
	}

	payload := make(map[string]string)
	payload["sub"] = string(idByte)
	payload["refreshTokenID"] = string(refreshTokenIDByte)
	tokenString, err := jwt.Signed(payload, expTime)
	if err != nil {
		logrus.Errorf("[GenerateAccessToken] signed error: %s", err)
		return nil, err
	}

	return &models.Token{
		AccessToken:      tokenString,
		RefreshToken:     refreshTokenString,
		ExpirationTime:   &expTime,
		ExpirationSecond: int(expTime.Sub(tn).Seconds()),
	}, nil
}

func (s *service) getExpireTime() time.Duration {
	return (time.Hour * (s.config.JWT.ExpireTime.Day * 24)) +
		(time.Hour * s.config.JWT.ExpireTime.Hour) +
		(time.Minute * s.config.JWT.ExpireTime.Minute)
}

// RefreshToken refresh token
func (s *service) RefreshToken(database *gorm.DB, refreshTokenString string) (*models.Token, error) {
	refreshToken := &models.RefreshToken{}
	err := s.tokenRepository.FindByRefreshTokenExpired(database, refreshTokenString, refreshToken)
	if err != nil {
		return nil, s.result.InvalidTokenNotExpire
	}

	token, err := s.GenerateAccessToken(database, refreshToken.UserID)
	if err != nil {
		logrus.Errorf("[RefreshToken] generate access token by userID=%d error: %s", refreshToken.UserID, err)
		return nil, err
	}

	refreshToken.Used = true
	err = s.tokenRepository.Delete(database, refreshToken)
	if err != nil {
		logrus.Errorf("[RefreshToken] remove refreshToken id=%d error: %s", refreshToken.ID, err)
		return nil, err
	}

	return token, nil
}
