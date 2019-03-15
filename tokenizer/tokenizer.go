package tokenizer

const (
	startVariable string = "$("
	endVariable   string = ")"
)

var str = ""

var (
	// Go doesn't allow
	startVarRunes = []rune(startVariable)
	endVarRunes   = []rune(endVariable)
)

func ParseResponse(response string) *ParsedResponse {
	parsed := &ParsedResponse{}
	parsed.RawText = response

	startIndex := 0
	currentIndex := 0
	runes := []rune(response)
	runeLen := len(runes)
	for currentIndex < runeLen {
		if isStartOfVariable(runes, currentIndex) {
			// Add previous strings to token slice.
			appendToken(parsed.Tokens, runes, startIndex, currentIndex)
			// Parse variable
			variableToken, endIndex := parseVariable(runes, currentIndex)
			parsed.Tokens = append(parsed.Tokens, variableToken)

			startIndex = endIndex
			currentIndex = endIndex
		} else {
			currentIndex++
		}
	}
	appendToken(parsed.Tokens, runes, startIndex, currentIndex)

	return parsed
}

func appendToken(tokens []Token, runes []rune, startIndex int, endIndex int) {
	if startIndex < endIndex {
		tokens = append(tokens, Token{
			RawText:   string(runes[startIndex:endIndex]),
			TokenType: TextTokenType,
		})
	}
}

// Returned int is start index of the new token (first index after the variable token)
// It is assumed that the runes from fromIndex starts with startVariable "$(".
func parseVariable(runes []rune, fromIndex int) (Token, int) {
	startIndex := len(startVarRunes)
	currentIndex := startIndex
	token := Token{Arguments: []Token{}}
	for currentIndex < len(runes) {
		if isStartOfVariable(runes, currentIndex) {
			appendToken(token.Arguments, runes, startIndex, currentIndex)
			subToken, nextIndex := parseVariable(runes, currentIndex)
			token.Arguments = append(token.Arguments, subToken)
			startIndex = nextIndex
			currentIndex = nextIndex
		} else {
			currentIndex++
		}
	}
	token.RawText = string(runes[startIndex:currentIndex])
	token.TokenType = VariableTokenType
	appendToken(token.Arguments, runes, startIndex, currentIndex)

	return token, currentIndex
}

func isStartOfVariable(runes []rune, fromIndex int) bool {
	// Check if there is no more room for variable
	// For example, If startVariable is "$(" and endVariable is ")",
	// At least 4 characters are needed for variable expression, like "$(a)"
	/*if len(runes) < fromIndex+len(startVarRunes)+len(endVarRunes)+1 {
		return false
	}*/
	return IsSubRune(runes, startVarRunes, fromIndex)
}

/*
Parse logic..

if $:
	check if next token is (
		If yes, recursively parse next tokens
	else [not ( or EOF]
		pass
else:
	append to the last token
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
