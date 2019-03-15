package tokenizer

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// Test parsing empty response
func TestNoVariableResponseEmpty(t *testing.T) {
	t.Parallel()
	response := ParseResponse("")
	assert.Equal(t, 1, len(response.Tokens))
	assert.Equal(t, Token{
		RawText:   "",
		TokenType: TextTokenType,
	}, response.Tokens[0])
}

// Test parsing one-word response, without variable
func TestNoVariableResponseOneWord(t *testing.T) {
	t.Parallel()
	response := ParseResponse("hello")
	assert.Equal(t, 1, len(response.Tokens))
	assert.Equal(t, Token{
		RawText:   "hello",
		TokenType: TextTokenType,
	}, response.Tokens[0])
}

// Test parsing response with multiple words, without variable
func TestNoVariableResponseWithSpaces(t *testing.T) {
	t.Parallel()
	response := ParseResponse("this is response text")
	assert.Equal(t, 1, len(response.Tokens))
	assert.Equal(t, Token{
		RawText:   "this is response text",
		TokenType: TextTokenType,
	}, response.Tokens[0])
}

// Test parsing response with special characters, without variable
func TestNoVariableResponseWithSpecialCharacters(t *testing.T) {
	t.Parallel()
	response := ParseResponse("!@#special $% ^&*()")
	assert.Equal(t, 1, len(response.Tokens))
	assert.Equal(t, Token{
		RawText:   "!@#special $% ^&*()",
		TokenType: TextTokenType,
	}, response.Tokens[0])
}

// response itself is a variable, "$(user)"
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
