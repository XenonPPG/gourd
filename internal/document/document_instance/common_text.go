package document_instance

import (
	"gourd/internal/common"
	"gourd/internal/document/rule"
	"gourd/internal/utils"
	"io"
)

type CommonTextDocument struct {
	content        string
	EscapeMarkdown bool
}

func (doc *CommonTextDocument) Read(in io.Reader) error {
	bytes, err := io.ReadAll(in)
	if err != nil {
		return err
	}
	doc.content = string(bytes)
	return nil
}

func (doc *CommonTextDocument) Apply(rules []rule.Rule) error {
	fileType := utils.Ternary(doc.EscapeMarkdown, common.Markdown, common.Text)

	for _, r := range rules {
		err := r.Apply(&doc.content, fileType)
		if err != nil {
			return err
		}
	}

	return nil
}

func (doc *CommonTextDocument) Write(out io.Writer) error {
	if sw, ok := out.(io.StringWriter); ok {
		_, err := sw.WriteString(doc.content)
		return err
	}
	_, err := out.Write([]byte(doc.content))
	return err
}
