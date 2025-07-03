package utils

import (
	"regexp"
	"strings"
)

func GenerateSeName(name string) string {
	return strings.ReplaceAll(strings.TrimSpace(name), " ", "_")
}

func GenerateSKU(name string) string {
	words := strings.Fields(name)
	sku := ""
	for _, w := range words {
		re := regexp.MustCompile(`[A-Za-z0-9]`)
		loc := re.FindString(w)
		if loc != "" {
			sku += strings.ToUpper(loc)
		}
	}
	return sku
}
