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

func markdownReplacer(templates []ReplacerTemplate) AnyAdapter {
	return Wrap(func(document *string) error {
		for _, template := range templates {
			*document = strings.Replace(
				*document,
				template.OldStr,
				utils.EscapeMarkdownV2(template.NewStr),
				-1,
			)
		}
		return nil
	})
}

func NewReplacer(templates []ReplacerTemplate) AdaptersMap {
	return AdaptersMap{
		domain.Docx:     docxReplacer(templates),
		domain.Markdown: markdownReplacer(templates),
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
