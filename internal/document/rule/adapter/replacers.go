package adapter

import (
	"gourd/internal/common"
	"sort"
	"strings"

	"github.com/fumiama/go-docx"
)

type ReplacerTemplate map[string]string

type replacePair struct {
	oldStr string
	newStr string
}

func makePairs(template ReplacerTemplate, escapeMarkdown bool) []replacePair {
	keys := make([]string, 0, len(template))
	for oldStr := range template {
		keys = append(keys, oldStr)
	}
	sort.Strings(keys)

	pairs := make([]replacePair, 0, len(keys))
	for _, oldStr := range keys {
		newStr := template[oldStr]
		if escapeMarkdown {
			newStr = common.EscapeMarkdownV2(newStr)
		}
		pairs = append(pairs, replacePair{oldStr: oldStr, newStr: newStr})
	}
	return pairs
}

func replaceAll(s string, pairs []replacePair) string {
	for _, p := range pairs {
		s = strings.ReplaceAll(s, p.oldStr, p.newStr)
	}
	return s
}

func docxReplacer(template ReplacerTemplate) AnyAdapter {
	pairs := makePairs(template, false)

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
					text.Text = replaceAll(text.Text, pairs)
				}
			}
		}
		return nil
	})
}

func commonTextReplacer(template ReplacerTemplate, escapeMarkdown bool) AnyAdapter {
	pairs := makePairs(template, escapeMarkdown)

	return Wrap(func(document *string) error {
		*document = replaceAll(*document, pairs)
		return nil
	})
}

func NewReplacer(template ReplacerTemplate) AdaptersMap {
	return AdaptersMap{
		common.Text:     commonTextReplacer(template, false),
		common.Markdown: commonTextReplacer(template, true),
		common.Docx:     docxReplacer(template),
	}
}

func NewSimpleReplacer(oldStr, newStr string) AdaptersMap {
	return NewReplacer(map[string]string{oldStr: newStr})
}
