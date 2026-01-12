package dashboard

import (
	"net/http"

	"github.com/direwen/go-server/internal/util"
	"github.com/labstack/echo/v4"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) GetDashboard(c echo.Context) error {
	result, err := h.service.GetPublicStats(c.Request().Context())
	if err != nil {
		return util.ErrorResponse(c, http.StatusInternalServerError, "an error occurred while fetching dashboard data", err)
	}

	return util.SuccessResponse(c, http.StatusOK, "dashboard data fetched successfully", result)
}
