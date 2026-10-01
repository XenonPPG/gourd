package rule

import (
	"errors"
	"fmt"
	"gourd/internal/document/rule/adapter"
	"gourd/internal/domain"
	"strings"
)

type Rule struct {
	ID               int
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
		Name: "Рекомендованное",
		Rules: []Rule{
			{
				Name:             "Заменить тире",
				Description:      "Заменяет — на –",
				Adapters:         adapter.NewSimpleReplacer("—", "–"),
				EnabledByDefault: true,
			},
			{
				Name:             "Заменить умножение",
				Description:      "Заменяет × на *",
				Adapters:         adapter.NewSimpleReplacer("×", "*"),
				EnabledByDefault: true,
			},
			{
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
				Name:        "Убрать ';'",
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

var Registry = assignIDs(categories)
var QuickRegistry = indexByID(Registry)

func assignIDs(cats []Category) []Category {
	nextRuleID := 1
	for i := range cats {
		for j := range cats[i].Rules {
			cats[i].Rules[j].ID = nextRuleID
			nextRuleID++
		}
	}
	return cats
}

func indexByID(cats []Category) map[int]Rule {
	reg := make(map[int]Rule)
	for _, cat := range cats {
		for _, r := range cat.Rules {
			reg[r.ID] = r
		}
	}
	return reg
}

func (r Rule) Apply(doc any, fileType domain.FileType) error {
	a, ok := r.Adapters[fileType]
	if !ok {
		return errors.New(fmt.Sprintf("adapter for %s not found", fileType))
	}
	return a.Apply(doc)
}
