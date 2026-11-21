package document

import (
	"gourd/internal/document/result"
	"gourd/internal/document/rule"
	"io"
)

type Document interface {
	Read(in io.Reader) error
	Write(out io.Writer) error
	Apply(rules []rule.Rule) ([]result.Result, error)
}
