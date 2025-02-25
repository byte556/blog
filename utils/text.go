package utils

import (
	"regexp"
)

// tagRe ищет любые HTML-теги вида <...>
var tagRe = regexp.MustCompile(`<[^>]+>`)

// StripHTMLTags убирает HTML-теги
func StripHTMLTags(s string) string {
	return tagRe.ReplaceAllString(s, "")
}

// Shorten обрезает строку до maxLen и добавляет "..."
func Shorten(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}
