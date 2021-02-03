package models

import (
	"time"
)

// Token token model
type Token struct {
	AccessToken      string     `json:"accessToken"`
	RefreshToken     string     `json:"refreshToken"`
	ExpirationTime   *time.Time `json:"expirationTime"`
	ExpirationSecond int        `json:"expirationSecond"`
}

// RefreshToken refreshToken model
type RefreshToken struct {
	Model
	RefreshToken string    `json:"refreshToken"`
	Used         bool      `json:"used"`
	ExpireTime   time.Time `json:"expireTime"`
	UserID       uint      `json:"userID"`
}
