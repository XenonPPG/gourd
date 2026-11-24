package dto

// Rule represents rule general information
type Rule struct {
	ID               string   `json:"id,omitempty"`
	Name             string   `json:"name,omitempty"`
	Description      string   `json:"description,omitempty"`
	EnabledByDefault bool     `json:"enabled_by_default,omitempty"`
	AvailableFormats []string `json:"available_formats,omitempty"`
}
