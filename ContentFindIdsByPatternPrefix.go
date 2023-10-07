package cms

import "regexp"

// returns the IDs in the content who have the following format [[prefix_id]]
func ContentFindIdsByPatternPrefix(content, prefix string) []string {
	ids := []string{}
	re := regexp.MustCompilePOSIX("|\\[\\[" + prefix + "_(.*)\\]\\]|U")
	matches := re.FindAllStringSubmatch(content, -1)
	for _, match := range matches {
		if match[0] == "" {
			continue
		}
		ids = append(ids, match[1])
	}

	return ids
}
