package parser

import (
	"strings"

	models "github.com/c-rainbow/simplechatbot/models"
)

/*

TODO: Arguments of variable is created differently from tokens of response,

When RawText has variable in it:
- Both response and variable creates tokens with texts and variables

However, when RawText does not have variable in it:
(1) Response creates a token with the entire plain text in it.
(2) Variable does not create a token with the entire plain text in it.


Should this be a problem?

*/

/*
To think about..
(1) nested variables $(followage $(user)), or problematic, $(urlfetch $(urlfetch $(urlfetch)...))
	conclusion: whitelist nestable variables. $user, $user_id, $display_name, $channel
(2) corrupt input, for example unclosed parentheses, "$(user]"
    conclusion: for v1, error.
(3) Unrecognized variable name, $(whatever)
    conclusion: for v1, error.
(4) Case-sensitive? $(user), $(User), $(uSeR)
    conclusion: for v1, case-sensitive.
(5) argument has start and end signs of variable (possibly with malicious reason), $(urlfetch "url_has_$(.html" )
    conclusion: for v1, error.
(6)
*/

const (
	VariableStartTag   string = "$("
	VariableClosingTag string = ")"
)

var (
	// Go doesn't allow creating slice in const section.
	startVarRunes = []rune(VariableStartTag)
	endVarRunes   = []rune(VariableClosingTag)
)

/*
Parse logic..

if current index is starting tag of a variable:
	Create currently unprocessed text as a token, add to token list.
	Parse the variable to a token, add to token list.
else:
	keep probing
finally:
	Create remaining unprocessed text as a token, add to token list.
*/
func ParseResponse(response string) *models.ParsedResponse {
	parsed := &models.ParsedResponse{}
	parsed.RawText = response

	// Working with multi-byte string in Go is painful. Let's just use runes.
	runes := []rune(response)
	startIndex := 0
	currentIndex := 0
	for currentIndex < len(runes) {
		if isVariableStartingTag(runes, currentIndex) {
			// Add previous strings to token slice.
			appendToken(&parsed.Tokens, runes, startIndex, currentIndex)

			// Parse variable
			variableToken, endIndex := parseVariable(runes, currentIndex)
			parsed.Tokens = append(parsed.Tokens, variableToken)

			// Adjust startIndex and currentIndex to the first index after variable
			startIndex = endIndex
			currentIndex = endIndex
		} else {
			currentIndex++
		}
	}
	// Create token for unprocessed runes, if any
	appendToken(&parsed.Tokens, runes, startIndex, currentIndex)

	return parsed
}

func appendToken(tokens *[]models.Token, runes []rune, startIndex int, endIndex int) {
	if startIndex < endIndex {
		*tokens = append(*tokens, models.Token{
			RawText:   string(runes[startIndex:endIndex]),
			TokenType: models.TextTokenType,
		})
	}
}

// Returned int is start index of the new token (first index after the variable token)
// It is assumed that the runes from fromIndex starts with startVariable "$(".
func parseVariable(runes []rune, fromIndex int) (models.Token, int) {
	startIndex := fromIndex + len(startVarRunes)
	currentIndex := startIndex

	token := models.Token{TokenType: models.VariableTokenType}
	// hasNestedVariable := false
	for currentIndex < len(runes) && !isVariableEndingTag(runes, currentIndex) {
		if isVariableStartingTag(runes, currentIndex) {
			// hasNestedVariable = true

			// Create token for previous runes
			appendToken(&token.Arguments, runes, startIndex, currentIndex)

			subToken, nextIndex := parseVariable(runes, currentIndex)
			token.Arguments = append(token.Arguments, subToken)

			startIndex = nextIndex
			currentIndex = nextIndex
		} else {
			currentIndex++
		}
	}

	// If there was no nested variable, then the for-loop above simply increments
	// currentIndex until it reached end of slice, or found ending tag for variable.
	// If there was nested variable, any texts after the last nested variable might
	// have been unprocessed as token.
	appendToken(&token.Arguments, runes, startIndex, currentIndex)

	// This increment may not need to be surrounded by if-statement, because
	// the only case when runes[currentIndex] is not end tag is end of string.
	if isVariableEndingTag(runes, currentIndex) {
		currentIndex += len(endVarRunes)
	}

	token.RawText = string(runes[fromIndex:currentIndex])
	token.VariableName = GetVariableName(&token)
	return token, currentIndex
}

// GetVariableName returns
func GetVariableName(token *models.Token) string {
	// Variable token with empty body "$()"
	if len(token.Arguments) == 0 {
		return ""
	}

	nameArg := token.Arguments[0]
	// The first argument of the variable token should always be text.
	// TokenType not being TextTokenType means nested variable in form of $($(a) b)
	if nameArg.TokenType != models.TextTokenType {
		return ""
	}

	// strings.Fields automatically deals with non-space whitespaces in the middle
	fields := strings.Fields(nameArg.RawText)
	if len(fields) == 0 {
		return ""
	}
	return fields[0]
}

func isVariableStartingTag(runes []rune, fromIndex int) bool {
	return IsSubRune(runes, startVarRunes, fromIndex)
}

func isVariableEndingTag(runes []rune, fromIndex int) bool {
	return IsSubRune(runes, endVarRunes, fromIndex)
}

// IsSubRune If runes is equal to subrunes from fromIndex.
func IsSubRune(runes []rune, subrunes []rune, fromIndex int) bool {
	if len(runes) < fromIndex+len(subrunes) {
		return false
	}
	// compare all indexes of runes starting from "fromIndex"
	for i, subRune := range subrunes {
		if runes[fromIndex+i] != subRune {
			return false
		}
	}
	return true
}
