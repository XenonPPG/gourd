package gateway

import (
	"gourd/internal/document/rule"
	"slices"

	"github.com/gofiber/fiber/v3"
)

type CategoryDTO struct {
	Name  string    `json:"name,omitempty"`
	Rules []RuleDTO `json:"rules,omitempty"`
}

type RuleDTO struct {
	ID               string   `json:"id,omitempty"`
	Name             string   `json:"name,omitempty"`
	Description      string   `json:"description,omitempty"`
	EnabledByDefault bool     `json:"enabled_by_default,omitempty"`
	AvailableFormats []string `json:"available_formats,omitempty"`
}

// ListRules returns processing rules grouped by category, optionally filtered by file type.
//
//	@Summary		List available rules
//	@Description	Returns the registry of document processing rules grouped by category.
//	@Description	If the optional `filter` query parameter is provided, only rules that have an adapter for that file type are included; otherwise all rules are returned.
//	@Description	Each category contains its name and a list of rules; each rule has an id, name, description, default enabled state and the list of file formats it is available for.
//	@Description	Categories without applicable rules are omitted from the response.
//	@Description	If no rules are found (the registry is empty or nothing matches the filter), the endpoint responds with 400.
//	@Tags			rules
//	@Produce		json
//	@Param			filter	query		string				false	"File type to filter rules by (only rules with an adapter for this type are returned)"
//	@Success		200		{array}		CategoryDTO			"List of rule categories with their rules"
//	@Failure		400		{object}	map[string]string	"No rules found for the given filter (response body: {\"error\": \"no rules found for file type\"})"
//	@Router			/rules [get]
func (s *Service) ListRules(c fiber.Ctx) error {
	filter := c.Query("filter")

	result := make([]CategoryDTO, 0, len(rule.Registry()))
	for _, rCat := range rule.Registry() {
		rules := make([]RuleDTO, 0, len(rCat.Rules))
		for _, r := range rCat.Rules {
			availableFileTypes := r.Adapters.AvailableAdapters()
			available := make([]string, 0, len(availableFileTypes))
			for _, t := range availableFileTypes {
				available = append(available, string(t))
			}

			if filter != "" && !slices.Contains(available, filter) {
				continue
			}

			rules = append(rules, RuleDTO{
				ID:               r.ID,
				Name:             r.Name,
				Description:      r.Description,
				EnabledByDefault: r.EnabledByDefault,
				AvailableFormats: available,
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
