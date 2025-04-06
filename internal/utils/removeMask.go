package utils

import "strings"

// RemoveMask remove caracteres não numéricos de uma string
func RemoveMask(value string) string {
	return strings.NewReplacer(".", "", "-", "", "/", "", "(", "", ")", "", " ", "").Replace(value)
}
