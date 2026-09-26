package document_instance

import (
	"gourd/internal/document/rule"
	"gourd/internal/domain"
	"io"
)

type MarkdownDocument struct {
	content string
}

func (doc *MarkdownDocument) Read(in io.Reader) error {
	bytes, err := io.ReadAll(in)
	if err != nil {
		return err
	}
	doc.content = string(bytes)
	return nil
}

func (doc *MarkdownDocument) Apply(rules []rule.Rule) error {
	for _, r := range rules {
		err := r.Apply(&doc.content, domain.Markdown)
		if err != nil {
			return err
		}
	}

	return nil
}

func (doc *MarkdownDocument) Write(out io.Writer) error {
	if sw, ok := out.(io.StringWriter); ok {
		_, err := sw.WriteString(doc.content)
		return err
	}
	_, err := out.Write([]byte(doc.content))
	return err
}
