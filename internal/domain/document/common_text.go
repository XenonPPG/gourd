package document

import (
	"gourd/internal/common"
	"gourd/internal/domain/result"
	"gourd/internal/domain/rule"
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

func (doc *CommonTextDocument) Apply(rules []rule.Rule) ([]result.Result, error) {
	fileType := utils.Ternary(doc.EscapeMarkdown, common.Markdown, common.Text)
	results := make([]result.Result, 0, len(rules))

	for _, r := range rules {
		res, err := r.Apply(&doc.content, fileType)
		if err != nil {
			return nil, err
		}
		if res != nil {
			results = append(results, result.Result{
				AnyResult: *res,
				RuleID:    r.ID,
			})
		}
	}

	return results, nil
}

func (doc *CommonTextDocument) Write(out io.Writer) error {
	if sw, ok := out.(io.StringWriter); ok {
		_, err := sw.WriteString(doc.content)
		return err
	}
	_, err := out.Write([]byte(doc.content))
	return err
}
