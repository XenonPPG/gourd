package domain

import (
	"gourd/internal/domain/result"
	"gourd/internal/domain/rule"
	"io"
)

type Document interface {
	Read(in io.Reader) error
	Write(out io.Writer) error
	Apply(rules []rule.Rule) ([]result.Result, error)
}
