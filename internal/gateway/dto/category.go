package dto

// Category represents Category, a group of rules (Rule specifically).
type Category struct {
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
	Rules       []Rule `json:"rules,omitempty"`
}
