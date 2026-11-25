package main

import (
	"fmt"
	"gourd/internal/common"
	"gourd/internal/domain/rule"
	"log"
	"os"
	"strings"
)

func main() {
	var text strings.Builder
	for _, cat := range rule.Registry() {
		// category name
		text.WriteString(fmt.Sprintf("### %s\n", cat.Name))

		// description
		text.WriteString(fmt.Sprintf("**%s**\n", cat.Description))

		// rules
		text.WriteString("|Правило|Эффект|\n|:-|-:|\n")
		for _, rule := range cat.Rules {
			text.WriteString(fmt.Sprintf(
				"|%s|%s|\n",
				rule.Name,
				common.EscapeMarkdownV2(rule.Description),
			))
		}
	}

	input, err := os.ReadFile("./cmd/makereadme/template.md")
	if err != nil {
		log.Fatal("Error reading file:", err)
	}

	output := strings.Replace(string(input), "{{RULES}}", text.String(), 1)

	err = os.WriteFile("./README.md", []byte(output), 0644)
	if err != nil {
		log.Fatal("Error writing file:", err)
	}
}
