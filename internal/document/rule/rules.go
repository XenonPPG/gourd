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
	Adapters         map[domain.FileType]adapter.AnyAdapter
	EnabledByDefault bool
}

var ruleDefaults = []struct {
	Enabled bool
	Rule    Rule
}{
	{true, newSimpleReplacer("Заменить тире", "Заменяет — на –", "—", "–")},
	{true, newSimpleReplacer("Заменить умножение", "Заменяет × на *", "×", "*")},
	{true, newReplacer("Заменить кавычки", "Заменяет «» на \"\"", []replacerTemplate{
		{oldStr: "«", newStr: "\""},
		{oldStr: "»", newStr: "\""},
	})},
	{true, newReplacer("Заменить стрелки", "Заменяет ← → на <- ->", []replacerTemplate{
		{oldStr: "←", newStr: "<-"},
		{oldStr: "→", newStr: "->"},
	})},
	{false, newLineEditor(
		"Убрать `;`",
		"Убирает точку с запятой в конце предложений",
		func(s string) (string, bool, error) {
			trimmed := strings.TrimSpace(s)
			if len(trimmed) == 0 {
				return "", false, nil
			}
			trimmed = strings.TrimSuffix(trimmed, ";")
			return trimmed, true, nil
		}),
	},
}

var Registry = buildRegistryMap()
var SortedRegistry = buildRegistry()

func buildRegistryMap() map[int]Rule {
	reg := make(map[int]Rule)
	for i, def := range ruleDefaults {
		def.Rule.EnabledByDefault = def.Enabled
		reg[i+1] = def.Rule
	}
	return reg
}

func buildRegistry() []Rule {
	reg := make([]Rule, 0, len(ruleDefaults))
	for i, def := range ruleDefaults {
		def.Rule.ID = i + 1
		def.Rule.EnabledByDefault = def.Enabled
		reg = append(reg, def.Rule)
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
