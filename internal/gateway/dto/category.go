package dto

// CategoryDTO represents Category, a group of rules (RuleDTO specifically).
type CategoryDTO struct {
	Name  string    `json:"name,omitempty"`
	Rules []RuleDTO `json:"rules,omitempty"`
}
