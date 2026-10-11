package adapter

import (
	"gourd/internal/common"
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

func commonTextReplacer(templates []ReplacerTemplate, escapeMarkdown bool) AnyAdapter {
	return Wrap(func(document *string) error {
		for _, template := range templates {
			var newStr string
			if escapeMarkdown {
				newStr = common.EscapeMarkdownV2(template.NewStr)
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
		common.Text:     commonTextReplacer(templates, false),
		common.Markdown: commonTextReplacer(templates, true),
		common.Docx:     docxReplacer(templates),
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
