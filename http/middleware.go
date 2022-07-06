package http

import (
	"mysql/app/apperr"
	"mysql/app/entity"

	"github.com/labstack/echo/v4"
)

const (
	claimsContextParam = "claims"
)

// RecoverPanicMiddleware is the middleware for handling panics.
func (s *ServerAPI) RecoverPanicMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {

		defer func() {
			if r := recover(); r != nil {
				err := apperr.Errorf(apperr.EUNKNOWN, "panic: %s", r)
				ErrorResponseJSON(c, err, nil)
			}
		}()

		return next(c)
	}
}

func AuthUser(c echo.Context) (*entity.User, error) {

	if claims, ok := c.Get(claimsContextParam).(*entity.AppClaims); ok {
		return claims.User, nil
	}

	// this should never happen
	return nil, apperr.Errorf(apperr.EUNAUTHORIZED, "no auth user found in context")
}
