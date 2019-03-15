package tokenizer

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
	startVariable string = "$("
	endVariable   string = ")"
)

var (
	// Go doesn't allow creating slice in const section.
	// These have to be in var section
	startVarRunes = []rune(startVariable)
	endVarRunes   = []rune(endVariable)
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
func ParseResponse(response string) *ParsedResponse {
	parsed := &ParsedResponse{}
	parsed.RawText = response

	startIndex := 0
	currentIndex := 0
	runes := []rune(response)
	for currentIndex < len(runes) {
		if isVariableStartTag(runes, currentIndex) {
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

func appendToken(tokens *[]Token, runes []rune, startIndex int, endIndex int) {
	if startIndex < endIndex {
		*tokens = append(*tokens, Token{
			RawText:   string(runes[startIndex:endIndex]),
			TokenType: TextTokenType,
		})
	}
}

// Returned int is start index of the new token (first index after the variable token)
// It is assumed that the runes from fromIndex starts with startVariable "$(".
func parseVariable(runes []rune, fromIndex int) (Token, int) {
	startIndex := fromIndex + len(startVarRunes)
	currentIndex := startIndex

	token := Token{TokenType: VariableTokenType}
	hasNestedVariable := false
	for currentIndex < len(runes) && !isVariableEndTag(runes, currentIndex) {
		if isVariableStartTag(runes, currentIndex) {
			hasNestedVariable = true

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
	if hasNestedVariable {
		appendToken(&token.Arguments, runes, startIndex, currentIndex)
	}

	// This if-statement might not be needed, because the only case when
	// runes[currentIndex] is not end tag is end of string.
	if isVariableEndTag(runes, currentIndex) {
		currentIndex++
	}

	token.RawText = string(runes[fromIndex:currentIndex])
	return token, currentIndex
}

func isVariableStartTag(runes []rune, fromIndex int) bool {
	return IsSubRune(runes, startVarRunes, fromIndex)
}

func isVariableEndTag(runes []rune, fromIndex int) bool {
	return IsSubRune(runes, endVarRunes, fromIndex)
}