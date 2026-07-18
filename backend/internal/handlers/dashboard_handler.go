package handlers

import (
	"sentineldesk/backend/internal/services"

	"github.com/gofiber/fiber/v2"
)

type DashboardHandler struct {
	service *services.DashboardService
}

func NewDashboardHandler() *DashboardHandler {
	return &DashboardHandler{
		service: services.NewDashboardService(),
	}
}

func (h *DashboardHandler) Summary(c *fiber.Ctx) error {

	summary, err := h.service.GetSummary()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(summary)
}

// MyIP returns the admin's IP and NetworkGroupID.
// network_group_id is determined by finding any device with the admin's IP
// and returning its network group. Frontend uses it for "Current Site" filtering.
// Falls back to empty string → All Devices mode.
func (h *DashboardHandler) MyIP(c *fiber.Ctx) error {
	ip := c.IP()
	groupID := h.service.GetAdminNetworkGroupID(ip)
	return c.JSON(fiber.Map{
		"ip":               ip,
		"network_group_id": groupID,
	})
}
