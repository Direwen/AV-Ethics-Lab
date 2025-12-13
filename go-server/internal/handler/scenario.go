package handler

import (
	"net/http"

	"github.com/direwen/go-server/internal/model"
	"github.com/direwen/go-server/internal/service"
	"github.com/direwen/go-server/internal/util"
	"github.com/labstack/echo/v4"
)

type ScenarioHandler struct {
	service service.ScenarioService
}

func NewScenarioHandler(service service.ScenarioService) *ScenarioHandler {
	return &ScenarioHandler{service: service}
}

func (h *ScenarioHandler) GetNext(c echo.Context) error {
	sessionID, _ := util.GetSessionID(c)
	scenario, err := h.service.GetNextScenario(
		c.Request().Context(),
		sessionID,
	)
	if err != nil {
		return err
	}

	if scenario == nil {
		return util.SuccessResponse(c, http.StatusOK, "no scenarios available", nil)
	}

	return util.SuccessResponse(c, http.StatusOK, "", map[string]*model.Scenario{"scenario": scenario})
}
