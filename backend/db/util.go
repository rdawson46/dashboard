package db

import (
	"regexp"
	ollama "github.com/ollama/ollama/api"
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

func generateDesc(message []ollama.Message) string {
	var lastQ string
	for _, m := range message {
		if m.Role == "user" {
			lastQ = m.Content
			break
		}
	}

	var desc string
	if len(lastQ) >= 10 {
		desc = lastQ[:10]
	} else {
		desc = lastQ
	}

	return desc
}
