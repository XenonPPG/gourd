package main

import (
	"fmt"
	"gourd/internal/document/rule"
	"gourd/internal/utils"
	"log"
	"os"
	"strings"
)

func main() {
	var text strings.Builder
	for _, cat := range rule.Registry {
		text.WriteString(fmt.Sprintf("### %s\n", cat.Name))
		for _, r := range cat.Rules {
			text.WriteString(fmt.Sprintf(
				"- %s: `%s`\n",
				utils.EscapeMarkdownV2(r.Name),
				utils.EscapeMarkdownV2(r.Description),
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
