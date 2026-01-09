package util

import (
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
)

type StandardResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

func SuccessResponse(c echo.Context, statusCode int, message string, data interface{}) error {
	return c.JSON(statusCode, StandardResponse{
		Success: true,
		Message: message,
		Data:    data,
	})
}

func ErrorResponse(c echo.Context, statusCode int, message string, err error) error {
	errorMsg := ""
	if err != nil {
		errorMsg = err.Error()
		// Log the error for debugging
		if statusCode >= 500 {
			log.Printf("Internal Server Error: %s - %v", message, err)
		}
	}
	return c.JSON(statusCode, StandardResponse{
		Success: false,
		Message: message,
		Error:   errorMsg,
	})
}

func CustomEchoErrorHandler(err error, c echo.Context) {
	// Default to 500 Internal Server Error
	code := http.StatusInternalServerError
	message := "Internal Server Error"

	// Check if it's an Echo HTTPError to extract the actual code and message
	if he, ok := err.(*echo.HTTPError); ok {
		code = he.Code
		if msg, ok := he.Message.(string); ok {
			message = msg
		}
	}

	ErrorResponse(c, code, message, err)
}
