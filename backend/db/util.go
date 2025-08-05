package db

import (
	"regexp"
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
