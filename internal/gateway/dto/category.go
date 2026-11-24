package dto

// Category represents Category, a group of rules (Rule specifically).
type Category struct {
	Name  string `json:"name,omitempty"`
	Rules []Rule `json:"rules,omitempty"`
}
