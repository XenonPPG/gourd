package result

import "gourd/internal/domain/rule/adapter"

type Result struct {
	adapter.AnyResult
	RuleID string
}
