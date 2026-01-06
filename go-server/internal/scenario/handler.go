package scenario

import (
	"net/http"

	"github.com/direwen/go-server/internal/util"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) GetNext(c echo.Context) error {
	// get session id from context (set by auth middleware)
	sessionID, ok := util.GetSessionID(c)
	if !ok {
		return util.ErrorResponse(c, http.StatusUnauthorized, "Missing session", nil)
	}

	id, err := uuid.Parse(sessionID)
	if err != nil {
		return util.ErrorResponse(c, http.StatusBadRequest, "Invalid session ID format", err)
	}

	scenario, err := h.service.GetNextScenario(c.Request().Context(), id)
	if err != nil {
		switch err.Error() {
		case "experiment completed":
			return util.ErrorResponse(c, http.StatusConflict, "Experiment completed", err)
		case "session expired":
			return util.ErrorResponse(c, http.StatusUnauthorized, "Session expired", err)
		default:
			return util.ErrorResponse(c, http.StatusInternalServerError, "Failed to get next scenario", err)
		}
	}

	return util.SuccessResponse(c, http.StatusOK, "Scenario retrieved", scenario)
}
