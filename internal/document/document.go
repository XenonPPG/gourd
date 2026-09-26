package document

import (
	"gourd/internal/document/rule"
	"io"
)

type Document interface {
	Read(in io.Reader) error
	Write(out io.Writer) error
	Apply(rules []rule.Rule) error
}
