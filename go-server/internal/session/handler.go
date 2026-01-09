package session

import (
	"net/http"

	"github.com/direwen/go-server/internal/util"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

var validate = validator.New()

type Handler struct {
	service Service
}

func NewHandler(s Service) *Handler {
	return &Handler{service: s}
}

func (h *Handler) Create(c echo.Context) error {
	var input CreateSessionInput

	if err := c.Bind(&input); err != nil {
		return util.ErrorResponse(c, http.StatusBadRequest, "Invalid request payload", err)
	}

	if err := validate.Struct(input); err != nil {
		return util.ErrorResponse(c, http.StatusBadRequest, "Validation failed", err)
	}

	token, err := h.service.RegisterSession(c.Request().Context(), input)
	if err != nil {
		return util.ErrorResponse(c, http.StatusInternalServerError, "Failed to create session", err)
	}

	return util.SuccessResponse(c, http.StatusCreated, "Session created successfully", map[string]string{"token": token})
}

func (h *Handler) GetSessionFeedback(c echo.Context) error {
	sessionIDStr, ok := util.GetSessionID(c)
	if !ok {
		return util.ErrorResponse(c, http.StatusUnauthorized, "Missing session", nil)
	}
	sessionID, err := uuid.Parse(sessionIDStr)
	if err != nil {
		return util.ErrorResponse(c, http.StatusBadRequest, "Invalid session ID format", err)
	}

	response, err := h.service.GetSessionFeedback(c.Request().Context(), sessionID)
	if err != nil {
		return util.ErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve session feedback", err)
	}

	return util.SuccessResponse(c, http.StatusOK, "Session feedback retrieved successfully", response)

}
