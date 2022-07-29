package handlers

import (
	"togo/config"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

var isLoggedIn = middleware.JWTWithConfig(middleware.JWTConfig{
	Claims:     &jwtCustomClaims{},
	SigningKey: []byte(config.Values.JWTSecret),
})

// private ...
func private(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		user := c.Get("user").(*jwt.Token)
		claims := user.Claims.(*jwtCustomClaims)

		c.Set("user_id", claims.UserID)

		return next(c)
	}
}
