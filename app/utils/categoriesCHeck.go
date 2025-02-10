package utils

import "strings"

func CategoriesCheck(cat []string) bool {
	cag := []string{"tech", "general", "sports", "education", "health"}

	validCategories := make(map[string]bool)
	for _, category := range cag {
		validCategories[strings.ToLower(category)] = true

	}

	exist := make(map[string]bool)

	for _, category := range cat {
		if !validCategories[strings.ToLower(category)] {
			return false
		}
		if exist[category] {
			return false
		}
		exist[category] = true
	}

	return true
}
