package adapter

import (
	"gourd/internal/domain"
	"gourd/internal/utils"
	"strings"

	"github.com/fumiama/go-docx"
)

type ReplacerTemplate struct {
	OldStr string
	NewStr string
}

func docxReplacer(templates []ReplacerTemplate) AnyAdapter {
	return Wrap(func(document *docx.Docx) error {
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
							template.OldStr,
							template.NewStr,
							-1,
						)
					}
				}
			}
		}
		return nil
	})
}

func textReplacer(templates []ReplacerTemplate, escapeMarkdown bool) AnyAdapter {
	return Wrap(func(document *string) error {
		for _, template := range templates {
			var newStr string
			if escapeMarkdown {
				newStr = utils.EscapeMarkdownV2(template.NewStr)
			} else {
				newStr = template.NewStr
			}

			*document = strings.Replace(
				*document,
				template.OldStr,
				newStr,
				-1,
			)
		}
		return nil
	})
}

func NewReplacer(templates []ReplacerTemplate) AdaptersMap {
	return AdaptersMap{
		domain.Text:     textReplacer(templates, false),
		domain.Markdown: textReplacer(templates, true),
		domain.Docx:     docxReplacer(templates),
	}
}

func NewSimpleReplacer(oldStr, newStr string) AdaptersMap {
	return NewReplacer([]ReplacerTemplate{
		{
			OldStr: oldStr,
			NewStr: newStr,
		},
	})
}
