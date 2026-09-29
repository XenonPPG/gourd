package gateway

import (
	"gourd/internal/document/rule"

	"github.com/gofiber/fiber/v3"
)

type RuleDTO struct {
	ID               int    `json:"id,omitempty"`
	Name             string `json:"name,omitempty"`
	Description      string `json:"description,omitempty"`
	EnabledByDefault bool   `json:"enabled_by_default,omitempty"`
}

// ListRules returns all available processing rules
//
//	@Summary List available rules
//	@Description Returns the full registry of document processing rules, including their id, name, description and default enabled state
//	@Tags rules
//	@Produce json
//	@Success 200 {array} RuleDTO "List of rules"
//	@Router /rules [get]
func (s *Service) ListRules(c fiber.Ctx) error {
	result := make([]RuleDTO, 0)
	for _, r := range rule.SortedRegistry {
		result = append(result, RuleDTO{
			ID:               r.ID,
			Name:             r.Name,
			Description:      r.Description,
			EnabledByDefault: r.EnabledByDefault,
		})
	}
	return c.JSON(result)
}
