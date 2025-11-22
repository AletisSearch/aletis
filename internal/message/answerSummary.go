package message

import (
	_ "embed"
	"strings"
)

//go:embed system-prompts/AnswerSummary.md
var SystemAnswerSummary string

type AnswerSummaryData struct {
	Title string
	URL   string
	MD    string
}

func AnswerSummaryPrompt(data *[]AnswerSummaryData, query string) string {
	var user strings.Builder
	user.WriteString("Context:\n\"\"\"\n")

	var context strings.Builder
	for _, d := range *data {
		context.WriteRune('[')
		context.WriteString(d.Title)
		context.WriteString(" - ")
		context.WriteString(d.URL)
		context.WriteString("]\n")
		context.WriteString(d.MD)
		context.WriteString("\n\n")
		user.WriteString(context.String())
		context.Reset()
	}
	user.WriteString("\n\"\"\"\n\n")
	user.WriteString("Query:\n---\n")
	user.WriteString(query)
	user.WriteString("\n---")
	return user.String()
}
