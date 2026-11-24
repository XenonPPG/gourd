package dto

// Result is the result produced by a single processing rule.
type Result struct {
	Result any    `json:"result"`
	Key    string `json:"key"`
	RuleID string `json:"rule_id"`
}

// IndicesResult result of rules that return indices of applied effect.
type IndicesResult struct {
	Indices    []int `json:"indices"`
	TotalCount int   `json:"total_count"`
}
