package dto

// RuleDTO represents rule general information
type RuleDTO struct {
	ID               string   `json:"id,omitempty"`
	Name             string   `json:"name,omitempty"`
	Description      string   `json:"description,omitempty"`
	EnabledByDefault bool     `json:"enabled_by_default,omitempty"`
	AvailableFormats []string `json:"available_formats,omitempty"`
}
