package healthAPI

import (
	"github.com/gofiber/fiber/v2"
)

// GetHealthHandler @Summary      Check if the service is running.
// @Description  Health check
// @Tags         health
// @Produce      plain
// @Success      200 {string} string "OK"
// @Router       /api/health [get]
func GetHealthHandler(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).SendString("OK")
}
