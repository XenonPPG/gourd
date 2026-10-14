package rule

import (
	"errors"
	"fmt"
	"gourd/internal/common"
	"gourd/internal/document/rule/adapter"
	"strings"
)

type Rule struct {
	ID               string
	Name             string
	Description      string
	Adapters         adapter.AdaptersMap
	EnabledByDefault bool
}

type Category struct {
	Name  string
	Rules []Rule
}

var categories = []Category{
	{
		Name: "Рекомендации",
		Rules: []Rule{
			{
				ID:          "replace-all-dashes",
				Name:        "Заменить все тире",
				Description: "Заменяет все тире на знак минуса",
				Adapters: adapter.NewReplacer(
					[]adapter.ReplacerTemplate{
						{OldStr: "—", NewStr: "-"},
						{OldStr: "–", NewStr: "-"},
					},
				),
				EnabledByDefault: true,
			},
			{
				ID:          "replace-en-dash",
				Name:        "Заменить тире",
				Description: "Заменяет — на –",
				Adapters:    adapter.NewSimpleReplacer("—", "–"),
			},
			{
				ID:               "replace-multiplication",
				Name:             "Заменить умножение",
				Description:      "Заменяет × на *",
				Adapters:         adapter.NewSimpleReplacer("×", "*"),
				EnabledByDefault: true,
			},
			{
				ID:          "replace-quotes",
				Name:        "Заменить кавычки",
				Description: "Заменяет «» на \"\"",
				Adapters: adapter.NewReplacer([]adapter.ReplacerTemplate{
					{OldStr: "«", NewStr: "\""},
					{OldStr: "»", NewStr: "\""},
				}),
				EnabledByDefault: true,
			},
		},
	},
	{
		Name: "Другое",
		Rules: []Rule{
			{
				ID:          "replace-arrows",
				Name:        "Заменить стрелки",
				Description: "Заменяет ← → на <- ->",
				Adapters: adapter.NewReplacer([]adapter.ReplacerTemplate{
					{OldStr: "←", NewStr: "<-"},
					{OldStr: "→", NewStr: "->"},
				}),
			},
		},
	},
	{
		Name: "Особое",
		Rules: []Rule{
			{
				ID:          "remove-semicolon",
				Name:        "Убрать точку с запятой",
				Description: "Убирает точку с запятой в конце предложений",
				Adapters: adapter.NewLineEditor(func(s string) (string, bool, error) {
					trimmed := strings.TrimSpace(s)
					if len(trimmed) == 0 {
						return "", false, nil
					}
					trimmed = strings.TrimSuffix(trimmed, ";")
					return trimmed, true, nil
				}),
			},
		},
	},
}

func Registry() []Category {
	return categories
}

var QuickRegistry = indexByID(Registry())

func indexByID(cats []Category) map[string]Rule {
	reg := make(map[string]Rule)
	for _, cat := range cats {
		for _, r := range cat.Rules {
			reg[r.ID] = r
		}
	}
	return reg
}

func (r Rule) Apply(doc any, fileType common.FileType) error {
	a, ok := r.Adapters[fileType]
	if !ok {
		return errors.New(fmt.Sprintf("adapter for %s not found", fileType))
	}
	return a.Apply(doc)
}
