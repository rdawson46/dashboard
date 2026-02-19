package db

import (
	"context"
	"fmt"
	"os"
	"regexp"

	ollama "github.com/ollama/ollama/api"
	"github.com/rdawson46/dashboard/api"
)

func checkPassword(pw string) bool {
	patternCheck := func() bool {
		patterns := []*regexp.Regexp{
			regexp.MustCompile("[0-9]"), // numberCheck 
			regexp.MustCompile("[!@#$%^&*(),.?:{}|<>]"), // specialCheck 
			regexp.MustCompile("[A-Z]"), // upperCheck 
		}

		for _, p := range patterns {
			if !p.MatchString(pw) { return false }
		}

		return true
	}()

	return len(pw) >= 8 && patternCheck
}


var summaryPrompt = `
# Instructions
Summarize the follow user query with a 3-5 description/title.
Only respond with the description and nothing else.
Use title case for the description.

## Examples
User Query: "Tell me about the queen of England"
Summary: Queen of England Info

User Query: "Write a http server in go"
Summary: HTTP Server in Go

User Query: "Give me a recipe for pizza"
Summary: Pizza Recipe


# User Query
%s

# Summary
`

func generateDesc(message []ollama.Message) string {
	var lastQ string
	for _, m := range message {
		if m.Role == "user" {
			lastQ = m.Content
			break
		}
	}

	url := os.Getenv("OLLAMA_URL")

	if url == "" {
		var desc string
		if len(lastQ) >= 10 {
			desc = lastQ[:10]
		} else {
			desc = lastQ
		}

		return desc
	}

	//
	ctx := context.Background()

    oc, err := api.NewOllamaClient(url, nil)

	if err != nil {
		var desc string
		if len(lastQ) >= 10 {
			desc = lastQ[:10]
		} else {
			desc = lastQ
		}

		return desc
	}

	query := fmt.Sprintf(summaryPrompt, lastQ)

	summary, err := oc.Chat(ctx, query)

	if err != nil {
		var desc string
		if len(lastQ) >= 10 {
			desc = lastQ[:10]
		} else {
			desc = lastQ
		}

		return desc
	}

	return summary
}
