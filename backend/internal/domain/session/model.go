package session

import "time"

type JWT struct {
	Value  string
	UserID uint
}

type RefreshToken struct {
	ID        uint      `gorm:"primaryKey"`
	UserID    uint      `gorm:"not null;index"`
	Value     string    `gorm:"not null;uniqueIndex"`
	CreatedAt time.Time `gorm:"not null"`
	ExpiresAt time.Time `gorm:"not null;index"`
}

type TokenData struct {
	JWT          *JWT
	RefreshToken *RefreshToken
}

func NewTokenData(jwt *JWT, refreshToken *RefreshToken) *TokenData {
	return &TokenData{
		JWT:          jwt,
		RefreshToken: refreshToken,
	}
}
