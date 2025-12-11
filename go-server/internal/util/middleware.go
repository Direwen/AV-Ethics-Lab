package util

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

func GetSessionID(c echo.Context) (string, bool) {
	token, ok := c.Get("session").(*jwt.Token)
	if !ok {
		return "", false
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", false
	}

	id, ok := claims["id"].(string)
	if !ok {
		return "", false
	}

	return id, true
}
