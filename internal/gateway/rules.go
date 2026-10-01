package gateway

import (
	"gourd/internal/document/rule"

	"github.com/gofiber/fiber/v3"
)

type CategoryDTO struct {
	Name  string    `json:"name,omitempty"`
	Rules []RuleDTO `json:"rules,omitempty"`
}

type RuleDTO struct {
	ID               int    `json:"id,omitempty"`
	Name             string `json:"name,omitempty"`
	Description      string `json:"description,omitempty"`
	EnabledByDefault bool   `json:"enabled_by_default,omitempty"`
}

// ListRules returns all available processing rules grouped by category
//
//	@Summary List available rules
//	@Description Returns the full registry of document processing rules grouped by category. Each category contains its id, name and a list of rules; each rule has an id, name, description and default enabled state
//	@Tags rules
//	@Produce json
//	@Success 200 {array} CategoryDTO "List of rule categories with their rules"
//	@Router /rules [get]
func (s *Service) ListRules(c fiber.Ctx) error {
	result := make([]CategoryDTO, 0, len(rule.Registry))
	for _, rCat := range rule.Registry {
		rules := make([]RuleDTO, 0, len(rCat.Rules))
		for _, r := range rCat.Rules {
			rules = append(rules, RuleDTO{
				ID:               r.ID,
				Name:             r.Name,
				Description:      r.Description,
				EnabledByDefault: r.EnabledByDefault,
			})
		}

		result = append(result, CategoryDTO{
			Name:  rCat.Name,
			Rules: rules,
		})
	}

	return c.JSON(result)
}
