package models

// UserSession user session
type UserSession struct {
	ID             uint `json:"id"`
	RefreshTokenID uint `json:"refreshTokenID"`
}
