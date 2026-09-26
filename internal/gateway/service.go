package gateway

import "github.com/gofiber/fiber/v3"

// Health checks the service status
//
//	@Summary Health check
//	@Description Returns a simple message confirming that the service is up and running
//	@Tags service
//	@Produce plain
//	@Success 200 {string} string "Healthy :)"
//	@Router /health [get]
func Health(c fiber.Ctx) error {
	return c.SendString("Healthy :)")
}
