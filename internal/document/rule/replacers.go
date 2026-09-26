package rule

import (
	"gourd/internal/document/rule/adapter"
	"gourd/internal/domain"
	"gourd/internal/utils"
	"strings"

	"github.com/fumiama/go-docx"
)

type replacerTemplate struct {
	oldStr string
	newStr string
}

func docxReplacer(templates []replacerTemplate) adapter.AnyAdapter {
	return adapter.Wrap(func(document *docx.Docx) error {
		for _, item := range document.Document.Body.Items {
			para, ok := item.(*docx.Paragraph)
			if !ok {
				continue
			}
			for _, child := range para.Children {
				run, ok := child.(*docx.Run)
				if !ok {
					continue
				}
				for _, runChildren := range run.Children {
					text, ok := runChildren.(*docx.Text)
					if !ok {
						continue
					}
					for _, template := range templates {
						text.Text = strings.Replace(
							text.Text,
							template.oldStr,
							template.newStr,
							-1,
						)
					}
				}
			}
		}
		return nil
	})
}

func markdownReplacer(templates []replacerTemplate) adapter.AnyAdapter {
	return adapter.Wrap(func(document *string) error {
		for _, template := range templates {
			*document = strings.Replace(
				*document,
				template.oldStr,
				utils.EscapeMarkdownV2(template.newStr),
				-1,
			)
		}
		return nil
	})
}

func newReplacer(name, description string, templates []replacerTemplate) Rule {
	return Rule{
		Name:        name,
		Description: description,
		Adapters: map[domain.FileType]adapter.AnyAdapter{
			domain.Docx:     docxReplacer(templates),
			domain.Markdown: markdownReplacer(templates),
		},
	}
}

func newSimpleReplacer(name, description, oldStr, newStr string) Rule {
	return newReplacer(name, description, []replacerTemplate{
		{
			oldStr: oldStr,
			newStr: newStr,
		},
	})
}
