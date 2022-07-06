package entity

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
)

const (
	AccessTokenExpiration  = 60 * time.Minute           // 1 hour
	RefreshTokenExpiration = 60 * 24 * 15 * time.Minute // 15 days
)

// AppClaims is a custom claims type for JWT
// It contains the information about the user and the standard claims
type AppClaims struct {
	jwt.StandardClaims
	User *User `json:"user"`
}

// NewAppClaims creates a new AppClaims
func NewAppClaims(user *User, expiresAfterMinutes time.Duration) *AppClaims {
	return &AppClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(expiresAfterMinutes).UTC().Unix(),
			NotBefore: time.Now().UTC().Unix(),
			Subject:   fmt.Sprint(user.ID),
			Id:        uuid.NewString(),
			IssuedAt:  time.Now().UTC().Unix(),
			Issuer:    "go-2022-stage",
			Audience:  "go-2022-stage-api",
		},
		User: user,
	}
}

// TokenPair is a struct that contains the tokens and the expiration time
type Token struct {
	AccessToken string `json:"access_token"`
}
