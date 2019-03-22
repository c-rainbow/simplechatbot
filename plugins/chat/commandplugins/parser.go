package commandplugins

import (
	"strings"
)

func GetTargetCommandNameAndResponse(text string) (string, string) {
	// TODO: This function does not acknowledge consecutive whitespaces in response text.
	// For example, if user types "!addcom !hello Welcome  \t  $(user)     here!", then
	// the response will be shortened to "Welcome $(user) here!", removing all long whitespaces
	// between words.
	fields := strings.Fields(text)

	switch len(fields) {
	case 1:
		return "", ""
	case 2:
		return fields[1], ""
	default:
		response := strings.Join(fields[2:], " ")
		return fields[1], response
	}
}
