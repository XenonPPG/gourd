package adapter

import (
	"gourd/internal/common"
	"slices"
	"strings"

	"github.com/fumiama/go-docx"
)

type effectFunc func(string) (output string, toBreak bool, err error)

func commonTextLineEditor(effect effectFunc) AnyAdapter {
	return Wrap(Indices, func(document *string) ([]int, error) {
		lines := strings.Split(*document, "\n")
		idx := make([]int, 0)
		for i, line := range lines {
			newLine, toBreak, err := effect(line)
			if err != nil {
				return nil, err
			}
			lines[i] = newLine
			idx = append(idx, i)
			if toBreak {
				break
			}
		}
		*document = strings.Join(lines, "\n")
		return idx, nil
	})
}

func docxLineEditor(effect effectFunc) AnyAdapter {
	return Wrap(Count, func(document *docx.Docx) (int, error) {
		replaceCount := 0
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
						return 0, err
					}
					text.Text = newText
					replaceCount++
					if toBreak {
						break
					}
				}
			}
		}
		return replaceCount, nil
	})
}

func NewLineEditor(effect effectFunc) AdaptersMap {
	return AdaptersMap{
		common.Text:     commonTextLineEditor(effect),
		common.Markdown: commonTextLineEditor(effect),
		common.Docx:     docxLineEditor(effect),
	}
}
