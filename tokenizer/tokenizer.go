package tokenizer

import "fmt"

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
			// fmt.Println("in isSTartFoVariable")
			// Add previous strings to token slice.
			appendToken(&parsed.Tokens, runes, startIndex, currentIndex)
			// Parse variable
			variableToken, endIndex := parseVariable(runes, currentIndex)
			parsed.Tokens = append(parsed.Tokens, variableToken)

			startIndex = endIndex
			currentIndex = endIndex
		} else {
			// fmt.Println("currentIndex incremented", currentIndex)
			currentIndex++
		}
	}
	// fmt.Println("startIndex:", startIndex, ", endIndex:", currentIndex)
	appendToken(&parsed.Tokens, runes, startIndex, currentIndex)

	return parsed
}

func appendToken(tokens *[]Token, runes []rune, startIndex int, endIndex int) {
	// fmt.Println("tokens before:", tokens)
	if startIndex < endIndex {
		*tokens = append(*tokens, Token{
			RawText:   string(runes[startIndex:endIndex]),
			TokenType: TextTokenType,
		})
	}
	// fmt.Println("tokens after:", tokens)
}

// Returned int is start index of the new token (first index after the variable token)
// It is assumed that the runes from fromIndex starts with startVariable "$(".
func parseVariable(runes []rune, fromIndex int) (Token, int) {
	startIndex := fromIndex + len(startVarRunes)
	fmt.Println("startIndex:", startIndex, ", fromIdnex:", fromIndex)
	currentIndex := startIndex
	token := Token{TokenType: VariableTokenType}
	hasNestedVariable := false
	for currentIndex < len(runes) && !isEndOfVariable(runes, currentIndex) {
		if isStartOfVariable(runes, currentIndex) {
			hasNestedVariable = true
			fmt.Println("Start of variable index: ", currentIndex)
			appendToken(&token.Arguments, runes, startIndex, currentIndex)
			subToken, nextIndex := parseVariable(runes, currentIndex)
			token.Arguments = append(token.Arguments, subToken)
			startIndex = nextIndex
			currentIndex = nextIndex
			fmt.Println("End of variable index: ", currentIndex)
		} else {
			currentIndex++
		}
	}
	if hasNestedVariable {
		appendToken(&token.Arguments, runes, startIndex, currentIndex)
	}

	if isEndOfVariable(runes, currentIndex) {
		currentIndex++
	}

	token.RawText = string(runes[fromIndex:currentIndex])
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

func isEndOfVariable(runes []rune, fromIndex int) bool {
	return IsSubRune(runes, endVarRunes, fromIndex)
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
