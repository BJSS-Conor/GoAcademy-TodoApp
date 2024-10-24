package stringUtils

import "strings"

func IsEmptyOrWhitespace(str string) bool {
	return strings.TrimSpace(str) == ""
}
