package dto

// ResultDTO is the result produced by a single processing rule.
type ResultDTO struct {
	Result any    `json:"result"`
	Key    string `json:"key"`
	RuleID string `json:"rule_id"`
}

// IndicesResultDTO result of rules that return indices of applied effect.
type IndicesResultDTO struct {
	Indices    []int `json:"indices"`
	TotalCount int   `json:"total_count"`
}
