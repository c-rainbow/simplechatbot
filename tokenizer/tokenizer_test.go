package tokenizer

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Test parsing empty response
// TODO: Should response have 0 or 1 token in case of empty string?
func TestNoVariableResponseEmpty(t *testing.T) {
	response := ParseResponse("")
	assert.Empty(t, response.Tokens)
}

// Test parsing one-word response, without variable
func TestNoVariableResponseOneWord(t *testing.T) {
	response := ParseResponse("hello")
	assert.Equal(t, 1, len(response.Tokens))
	assert.Equal(t, Token{
		RawText:   "hello",
		TokenType: TextTokenType,
	}, response.Tokens[0])
}

// Test parsing response with multiple words, without variable
func TestNoVariableResponseWithSpaces(t *testing.T) {
	response := ParseResponse("this is response text")
	assert.Equal(t, 1, len(response.Tokens))
	assert.Equal(t, Token{
		RawText:   "this is response text",
		TokenType: TextTokenType,
	}, response.Tokens[0])
}

// Test parsing response with special characters, without variable
func TestNoVariableResponseWithSpecialCharacters(t *testing.T) {
	response := ParseResponse("!@#special $% ^&*()")
	assert.Equal(t, 1, len(response.Tokens))
	assert.Equal(t, Token{
		RawText:   "!@#special $% ^&*()",
		TokenType: TextTokenType,
	}, response.Tokens[0])
}

// response itself is a variable, "$(user)"
func TestVariableEntireResponse(t *testing.T) {
	response := ParseResponse("$(user)")
	assert.Equal(t, 1, len(response.Tokens))
	assert.Equal(t, Token{
		RawText:   "$(user)",
		TokenType: VariableTokenType,
	}, response.Tokens[0])
}

func TestVariablePartialResponse(t *testing.T) {
	response := ParseResponse("hi $(user) hello")
	assert.Equal(t, 3, len(response.Tokens))

	assert.Equal(t, Token{
		RawText:   "hi ",
		TokenType: TextTokenType,
	}, response.Tokens[0])

	assert.Equal(t, Token{
		RawText:   "$(user)",
		TokenType: VariableTokenType,
	}, response.Tokens[1])

	assert.Equal(t, Token{
		RawText:   " hello",
		TokenType: TextTokenType,
	}, response.Tokens[2])
}

// Test when the same variable is used multiple times
func TestVariableSameVariableMultipleTimes(t *testing.T) {
	response := ParseResponse("hi $(user) hello $(user)")
	assert.Equal(t, 4, len(response.Tokens))

	assert.Equal(t, Token{
		RawText:   "hi ",
		TokenType: TextTokenType,
	}, response.Tokens[0])

	assert.Equal(t, Token{
		RawText:   "$(user)",
		TokenType: VariableTokenType,
	}, response.Tokens[1])

	assert.Equal(t, Token{
		RawText:   " hello ",
		TokenType: TextTokenType,
	}, response.Tokens[2])

	assert.Equal(t, Token{
		RawText:   "$(user)",
		TokenType: VariableTokenType,
	}, response.Tokens[3])
}

func TestVariableMultipleVariables(t *testing.T) {
	response := ParseResponse("hi $(user) to $(channel)")
	assert.Equal(t, 4, len(response.Tokens))

	assert.Equal(t, Token{
		RawText:   "hi ",
		TokenType: TextTokenType,
	}, response.Tokens[0])

	assert.Equal(t, Token{
		RawText:   "$(user)",
		TokenType: VariableTokenType,
	}, response.Tokens[1])

	assert.Equal(t, Token{
		RawText:   " to ",
		TokenType: TextTokenType,
	}, response.Tokens[2])

	assert.Equal(t, Token{
		RawText:   "$(channel)",
		TokenType: VariableTokenType,
	}, response.Tokens[3])
}

func TestNestedVariable(t *testing.T) {
	response := ParseResponse("hi $(urlfetch http://twitch.tv/$(user)/01)")
	assert.Equal(t, 2, len(response.Tokens))

	assert.Equal(t, Token{
		RawText:   "hi ",
		TokenType: TextTokenType,
	}, response.Tokens[0])

	vToken := response.Tokens[1]
	assert.Equal(t, "$(urlfetch http://twitch.tv/$(user)/01)", vToken.RawText)
	assert.Equal(t, VariableTokenType, vToken.TokenType)
	fmt.Println("vToken: ", vToken)
	assert.Equal(t, 3, len(vToken.Arguments))

	// Examine nested tokens
	assert.Equal(t, Token{
		RawText:   "urlfetch http://twitch.tv/",
		TokenType: TextTokenType},
		vToken.Arguments[0])
	assert.Equal(t, Token{
		RawText:   "$(user)",
		TokenType: VariableTokenType},
		vToken.Arguments[1])
	assert.Equal(t, Token{
		RawText:   "/01",
		TokenType: TextTokenType},
		vToken.Arguments[2])

}

// response with one variable, "welcome $(user)", "$(user) hi"
// response with multiple variables of same type, "Display name of $(userid) is $(display_name)"
// response with multiple variables of different types, "$(user) has followed for $(follow_age)"
// nested variables, "User has followed for $(follow_age $(user))""
// Mix of texts and variables without space "A$(user)B$(channel)C"
// Continuous variables "$(display_name)$(user)$(follow_age $(user))"
// Ending with variable
// Variable error cases..
// one-character variable name  "hello$(n)"

// Test isStartOfVariable(), in various cases
