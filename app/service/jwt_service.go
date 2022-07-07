package service

import (
	"context"
	"time"

	"mysql/app/entity"
)

// JWTService is an interface for JWT service.
// It is used to generate and validate JWT tokens.
type JWTService interface {
	// Exchange a auth entity for a JWT token pair.
	Exchange(ctx context.Context, auth *entity.User) (*entity.Token, error)

	// Parse a JWT token and return the associated claims.
	Parse(ctx context.Context, token string) (*entity.AppClaims, error)
}

// JWTBlacklistService is an interface for JWT blacklist service.
type JWTBlacklistService interface {

	// Invalidate a JWT token.
	// Returns EUNAUTHORIZED if the user is not allowed to invalidate the token.
	Invalidate(ctx context.Context, token string, expiration time.Duration) error

	// IsBlacklisted checks if a token is blacklisted.
	IsBlacklisted(ctx context.Context, token string) (bool, error)
}