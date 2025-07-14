package utils

import (
	"regexp"
	"strconv"
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

func GenerateUniqueSeName(baseSeName string, existingMatches map[string]struct{}) string {
	seName := baseSeName
	i := 1
	for {
		if _, exists := existingMatches[seName]; !exists {
			break
		}
		seName = baseSeName + "_" + strconv.Itoa(i)
		i++
	}

	return seName
}

func GenerateUniqueSku(baseSku string, existingMatches map[string]struct{}) string {
	sku := baseSku
	i := 1
	for {
		if _, exists := existingMatches[sku]; !exists {
			break
		}
		sku = baseSku + "-" + strconv.Itoa(i)
		i++
	}

	return sku
}
