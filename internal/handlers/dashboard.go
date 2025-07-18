package handlers

import (
	"github.com/ekastn/hms-api/internal/service"
	"github.com/ekastn/hms-api/internal/utils"
	"github.com/gofiber/fiber/v2"
)

type DashboardHandler struct {
	dashboardService *service.DashboardService
}

func NewDashboardHandler(dashboardService *service.DashboardService) *DashboardHandler {
	return &DashboardHandler{
		dashboardService: dashboardService,
	}
}

// GetDashboardData handles the request to get dashboard data.
//
//	@Summary		Get dashboard data
//	@Description	Retrieve various statistics and recent activities for the dashboard. Admin or Management access required.
//	@Tags			Dashboard
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Success		200	{object}	domain.DashboardResponse	"Dashboard data retrieved successfully"
//	@Failure		500	{object}	utils.ErrorResponse			"Failed to get dashboard data"
//	@Router			/dashboard [get]
func (h *DashboardHandler) GetDashboardData(c *fiber.Ctx) error {
	data, err := h.dashboardService.GetDashboardData(c.Context())
	if err != nil {
		return utils.ErrorResponseJSON(c, fiber.StatusInternalServerError, "Failed to get dashboard data", nil)
	}

	return c.JSON(data)
}
