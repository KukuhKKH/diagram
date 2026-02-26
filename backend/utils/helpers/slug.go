package helpers

import (
	"regexp"
	"strings"
)

var (
	regexpInvalidChars = regexp.MustCompile(`[^a-z0-9\s-]`)
	regexpSpaces       = regexp.MustCompile(`\s+`)
)

func Slug(s string) string {
	s = strings.ToLower(s)
	s = regexpInvalidChars.ReplaceAllString(s, "")
	s = regexpSpaces.ReplaceAllString(s, "-")
	s = strings.Trim(s, "-")

	return s
}
