package gateway

import (
	"gourd/internal/document/rule"
	"gourd/internal/domain"

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

// ListRules returns processing rules available for the given file type, grouped by category.
//
//	@Summary		List available rules for a file type
//	@Description	Returns the registry of document processing rules grouped by category, filtered by the file type from the path. Each category contains its name and a list of rules; each rule has an id, name, description and default enabled state. Only rules that have an adapter for the requested file type are included. Categories without applicable rules are still returned, with an empty rules list.
//	@Tags			rules
//	@Produce		json
//	@Param			type	path		string					true	"File type the rules are requested for"
//	@Success		200		{array}		CategoryDTO				"List of rule categories with their rules"
//	@Failure		404		{object}	map[string]string		"Rule registry contains no categories (response body: {\"error\": \"no rules found for file type\"})"
//	@Router			/rules/{type} [get]
func (s *Service) ListRules(c fiber.Ctx) error {
	fileType := domain.FileType(c.Params("type"))

	result := make([]CategoryDTO, 0, len(rule.Registry))
	for _, rCat := range rule.Registry {
		rules := make([]RuleDTO, 0, len(rCat.Rules))
		for _, r := range rCat.Rules {
			_, ok := r.Adapters[fileType]
			if !ok {
				continue
			}

			rules = append(rules, RuleDTO{
				ID:               r.ID,
				Name:             r.Name,
				Description:      r.Description,
				EnabledByDefault: r.EnabledByDefault,
			})
		}
		if len(rules) == 0 {
			continue
		}

		result = append(result, CategoryDTO{
			Name:  rCat.Name,
			Rules: rules,
		})
	}

	if len(result) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "no rules found for file type",
		})
	}

	return c.JSON(result)
}
