package cms

import "regexp"

func isNumeric(str string) bool {

	var digitCheck = regexp.MustCompile(`^[0-9]+$`)

	return digitCheck.MatchString(str)
}
