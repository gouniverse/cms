package cms

import (
	"regexp"
	"strings"
)

func isNumeric(str string) bool {

	var digitCheck = regexp.MustCompile(`^[0-9]+$`)

	return digitCheck.MatchString(str)
}

// isJSON is naive implementation for superficial, rough and fast checking for JSON
func isJSON(str string) bool {
	if strings.HasPrefix(str, "{") && strings.HasSuffix(str, "}") {
		return true
	}

	if strings.HasPrefix(str, "[") && strings.HasSuffix(str, "]") {
		return true
	}

	return false
}
