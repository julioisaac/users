package healthAPI

import (
	"github.com/gofiber/fiber/v2"
)

const (
	HandlerPath = "/health"
)

func RegisterRoutes(router fiber.Router) {
	router.Get("", GetHealthHandler)
}
