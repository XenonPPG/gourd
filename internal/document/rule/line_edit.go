package rule

import (
	"gourd/internal/document/rule/adapter"
	"gourd/internal/domain"
	"slices"
	"strings"

	"github.com/fumiama/go-docx"
)

type effectFunc func(string) (output string, toBreak bool, err error)

func markdownLineEditor(effect effectFunc) adapter.AnyAdapter {
	return adapter.Wrap(func(document *string) error {
		lines := strings.Split(*document, "\n")
		for i, line := range lines {
			newLine, toBreak, err := effect(line)
			if err != nil {
				return err
			}
			lines[i] = newLine
			if toBreak {
				break
			}
		}
		*document = strings.Join(lines, "\n")
		return nil
	})
}

func docxLineEditor(effect effectFunc) adapter.AnyAdapter {
	return adapter.Wrap(func(document *docx.Docx) error {
		for _, item := range document.Document.Body.Items {
			para, ok := item.(*docx.Paragraph)
			if !ok {
				continue
			}
			// reverse order helps specific rule and doesn't affect other
			for _, v := range slices.Backward(para.Children) {
				run, ok := v.(*docx.Run)
				if !ok {
					continue
				}
				for _, v := range slices.Backward(run.Children) {
					text, ok := v.(*docx.Text)
					if !ok {
						continue
					}
					newText, toBreak, err := effect(text.Text)
					if err != nil {
						return err
					}
					text.Text = newText
					if toBreak {
						break
					}
				}
			}
		}
		return nil
	})
}

func newLineEditor(name, description string, effect effectFunc) Rule {
	return Rule{
		Name:        name,
		Description: description,
		Adapters: map[domain.FileType]adapter.AnyAdapter{
			domain.Markdown: markdownLineEditor(effect),
			domain.Docx:     docxLineEditor(effect),
		},
	}
}
