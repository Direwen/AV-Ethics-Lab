package response

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

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) Create(c echo.Context) error {
	// get session id from context (set by auth middleware)
	sessionIDStr, ok := util.GetSessionID(c)
	if !ok {
		return util.ErrorResponse(c, http.StatusUnauthorized, "Missing session", nil)
	}
	sessionID, err := uuid.Parse(sessionIDStr)
	if err != nil {
		return util.ErrorResponse(c, http.StatusBadRequest, "Invalid session ID format", err)
	}

	scenarioIDStr := c.Param("scenario_id")
	scenarioID, err := uuid.Parse(scenarioIDStr)
	if err != nil {
		return util.ErrorResponse(c, http.StatusBadRequest, "Invalid scenario ID format", err)
	}

	var input SubmitResponseInput
	if err := c.Bind(&input); err != nil {
		return util.ErrorResponse(c, http.StatusBadRequest, "Invalid input", err)
	}

	if err := validate.Struct(input); err != nil {
		return util.ErrorResponse(c, http.StatusBadRequest, "Invalid input", err)
	}

	response, err := h.service.SubmitResponse(c.Request().Context(), sessionID, scenarioID, input)
	if err != nil {
		return util.ErrorResponse(c, http.StatusBadRequest, "Failed to submit response", err)
	}

	return util.SuccessResponse(c, http.StatusCreated, "Response submitted", response)
}
