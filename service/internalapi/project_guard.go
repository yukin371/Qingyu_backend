package internalapi

import "strings"

func sameProjectID(actual, expected string) bool {
	return strings.TrimSpace(actual) == strings.TrimSpace(expected)
}
