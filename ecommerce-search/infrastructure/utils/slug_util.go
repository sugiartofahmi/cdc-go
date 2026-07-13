package utils

import "strings"

func GenerateSlug(text string) string {
	slug := strings.ToLower(text)
	slug = strings.TrimSpace(slug)
	slug = strings.ReplaceAll(slug, " ", "-")
	return slug
}
