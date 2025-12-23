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
	sessionID, _ := util.GetSessionID(c)
	s, err := h.service.GetNextScenario(
		c.Request().Context(),
		uuid.MustParse(sessionID),
	)
	if err != nil {
		return err
	}

	if s == nil {
		return util.SuccessResponse(c, http.StatusOK, "no scenarios available", nil)
	}

	return util.SuccessResponse(c, http.StatusOK, "", map[string]*Scenario{"scenario": s})
}
