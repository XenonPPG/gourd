package result

import "gourd/internal/document/rule/adapter"

type Result struct {
	adapter.AnyResult
	RuleID string
}
