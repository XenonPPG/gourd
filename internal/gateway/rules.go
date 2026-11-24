package gateway

import (
	"gourd/internal/domain/rule"

	"github.com/gofiber/fiber/v3"
)

// ListRules returns all processing rules grouped by category.
//
//	@Summary		List available rules
//	@Description	Returns the registry of domain processing rules grouped by category.
//	@Description	Each category contains its name and a list of rules; each rule has an id, name, description, default enabled state and the list of file formats it is available for.
//	@Description	The categories are returned under the `categories` key of the response object.
//	@Description	If the rules cannot be built from the registry, the endpoint responds with 500.
//	@Tags			rules
//	@Produce		json
//	@Success		200	{object}	map[string][]dto.Category	"Rule categories with their rules (response body: {\"categories\": [...]})"
//	@Failure		500	{object}	map[string]string				"Failed to build rules from the registry (response body: {\"error\": \"...\"})"
//	@Router			/rules [get]
func (s *Service) ListRules(c fiber.Ctx) error {
	cats, err := s.dtoRegistry.Match(rule.Registry())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err,
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"categories": cats,
	})
}
