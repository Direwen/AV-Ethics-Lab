package middleware

import (
	"log"
	"net/http"
	"os"

	"github.com/direwen/go-server/internal/util"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
)

func JWTMiddleware() echo.MiddlewareFunc {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		log.Fatal("JWT_SECRET is not set")
	}

	config := echojwt.Config{
		SigningKey: []byte(secret),
		ContextKey: "session",
		ErrorHandler: func(c echo.Context, err error) error {
			return util.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized", err)
		},
	}

	return echojwt.WithConfig(config)
}
