package main

import (
	"fmt"
	"gourd/internal/document/rule"
	"gourd/internal/utils"
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
		fmt.Println("Error reading file:", err)
		return
	}

	output := strings.Replace(string(input), "{{RULES}}", text.String(), 1)

	err = os.WriteFile("./README.md", []byte(output), 0644)
	if err != nil {
		fmt.Println("Error writing file:", err)
		return
	}
}
