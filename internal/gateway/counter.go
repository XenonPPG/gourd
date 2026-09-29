package gateway

import "github.com/gofiber/fiber/v3"

// GetCounterValue returns the current counter value
//
//	@Summary Get counter value
//	@Description Returns the current value of the counter
//	@Tags counter
//	@Produce json
//	@Success 200 {object} int "Current counter value"
//	@Router /counter [get]
func (s *Service) GetCounterValue(c fiber.Ctx) error {
	return c.JSON(s.counterService.GetValue())
}
