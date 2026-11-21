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

func replaceIdx(s, old, repl string) (string, []int) {
	if old == "" {
		return s, nil
	}
	var idx []int
	for pos := 0; ; {
		i := strings.Index(s[pos:], old)
		if i < 0 {
			break
		}
		idx = append(idx, pos+i)
		pos += i + len(old)
	}
	if len(idx) == 0 {
		return s, nil
	}

	var b strings.Builder
	b.Grow(len(s) + len(idx)*(len(repl)-len(old)))
	prev := 0
	for _, i := range idx {
		b.WriteString(s[prev:i])
		b.WriteString(repl)
		prev = i + len(old)
	}
	b.WriteString(s[prev:])

	return b.String(), idx
}

func replaceAllIdx(s string, pairs []replacePair) (string, []int) {
	idx := make([]int, 0)
	i := make([]int, 0)
	for _, p := range pairs {
		s, i = replaceIdx(s, p.oldStr, p.newStr)
		idx = append(idx, i...)
	}
	return s, idx
}

func docxReplacer(template ReplacerTemplate) AnyAdapter {
	pairs := makePairs(template, false)

	return Wrap(Count, func(document *docx.Docx) (int, error) {
		counter := 0

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
					idx := make([]int, 0)
					text.Text, idx = replaceAllIdx(text.Text, pairs)
					counter += len(idx)
				}
			}
		}
		return counter, nil
	})
}

func commonTextReplacer(template ReplacerTemplate, escapeMarkdown bool) AnyAdapter {
	pairs := makePairs(template, escapeMarkdown)

	return Wrap(Indices, func(document *string) ([]int, error) {
		idx := make([]int, 0)
		*document, idx = replaceAllIdx(*document, pairs)
		return idx, nil
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
