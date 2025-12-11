package handler

import (
	"net/http"

	"github.com/direwen/go-server/internal/dto"
	"github.com/direwen/go-server/internal/service"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

var validate = validator.New()

type SessionHandler struct {
	service service.SessionService
}

func NewSessionHandler(s service.SessionService) *SessionHandler {
	return &SessionHandler{service: s}
}

func (h *SessionHandler) Create(c echo.Context) error {
	var input dto.CreateSessionInput

	// Bind request data to input
	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request payload"})
	}

	// Validate struct defined fields with validate tags
	if err := validate.Struct(input); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	// register the session
	token, err := h.service.RegisterSession(c.Request().Context(), input)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, map[string]string{"token": token})
}
