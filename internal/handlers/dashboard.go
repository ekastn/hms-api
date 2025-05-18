package handlers

import (
	"github.com/ekastn/hms-api/internal/service"
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

func (h *DashboardHandler) GetDashboardData(c *fiber.Ctx) error {
	data, err := h.dashboardService.GetDashboardData(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get dashboard data: " + err.Error(),
		})
	}

	return c.JSON(data)
}
