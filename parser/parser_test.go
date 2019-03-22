package parser

import (
	"testing"

	models "github.com/c-rainbow/simplechatbot/models"
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

	assert.Equal(t, "hello", response.RawText)
	assert.Equal(t, 1, len(response.Tokens))
	assert.Equal(t, models.Token{
		RawText:   "hello",
		TokenType: models.TextTokenType,
	}, response.Tokens[0])
}

// Test parsing response with multiple words, without variable
func TestNoVariableResponseWithSpaces(t *testing.T) {
	response := ParseResponse("this is response text")

	assert.Equal(t, "this is response text", response.RawText)
	assert.Equal(t, 1, len(response.Tokens))
	assert.Equal(t, models.Token{
		RawText:   "this is response text",
		TokenType: models.TextTokenType,
	}, response.Tokens[0])
}

// Test parsing response with special characters, without variable
func TestNoVariableResponseWithSpecialCharacters(t *testing.T) {
	response := ParseResponse("!@#special $% ^&*()")

	assert.Equal(t, "!@#special $% ^&*()", response.RawText)
	assert.Equal(t, 1, len(response.Tokens))
	assert.Equal(t, models.Token{
		RawText:   "!@#special $% ^&*()",
		TokenType: models.TextTokenType,
	}, response.Tokens[0])
}

// response itself is a variable, "$(user)"
func TestVariableEntireResponse(t *testing.T) {
	response := ParseResponse("$(user)")

	assert.Equal(t, "$(user)", response.RawText)
	assert.Equal(t, 1, len(response.Tokens))
	assert.Equal(t, models.Token{
		RawText:      "$(user)",
		TokenType:    models.VariableTokenType,
		VariableName: "user",
		Arguments: []models.Token{
			models.Token{
				RawText:   "user",
				TokenType: models.TextTokenType,
			},
		},
	}, response.Tokens[0])
}

func TestVariableEntireResponseWithSpaceInName(t *testing.T) {
	response := ParseResponse("$(  user   )")

	assert.Equal(t, "$(  user   )", response.RawText)
	assert.Equal(t, 1, len(response.Tokens))
	assert.Equal(t, models.Token{
		RawText:      "$(  user   )",
		TokenType:    models.VariableTokenType,
		VariableName: "user",
		Arguments: []models.Token{
			models.Token{
				RawText:   "  user   ",
				TokenType: models.TextTokenType,
			},
		},
	}, response.Tokens[0])
}

func TestVariablePartialResponse(t *testing.T) {
	response := ParseResponse("hi $(user) hello")

	assert.Equal(t, "hi $(user) hello", response.RawText)
	assert.Equal(t, 3, len(response.Tokens))

	assert.Equal(t, models.Token{
		RawText:   "hi ",
		TokenType: models.TextTokenType,
	}, response.Tokens[0])

	assert.Equal(t, models.Token{
		RawText:      "$(user)",
		TokenType:    models.VariableTokenType,
		VariableName: "user",
		Arguments: []models.Token{
			models.Token{
				RawText:   "user",
				TokenType: models.TextTokenType,
			},
		},
	}, response.Tokens[1])

	assert.Equal(t, models.Token{
		RawText:   " hello",
		TokenType: models.TextTokenType,
	}, response.Tokens[2])
}

// one-character variable name  "hello$(n)"
func TestVariableShortName(t *testing.T) {
	response := ParseResponse("hello$(n)")

	assert.Equal(t, "hello$(n)", response.RawText)
	assert.Equal(t, 2, len(response.Tokens))

	assert.Equal(t, models.Token{
		RawText:   "hello",
		TokenType: models.TextTokenType,
	}, response.Tokens[0])

	assert.Equal(t, models.Token{
		RawText:      "$(n)",
		TokenType:    models.VariableTokenType,
		VariableName: "n",
		Arguments: []models.Token{
			models.Token{
				RawText:   "n",
				TokenType: models.TextTokenType,
			},
		},
	}, response.Tokens[1])
}

// Test when the same variable is used multiple times
// "welcome $(user)", "$(user) hi"
func TestVariableSameVariableMultipleTimes(t *testing.T) {
	response := ParseResponse("hi $(user) hello $(user)")

	assert.Equal(t, "hi $(user) hello $(user)", response.RawText)
	assert.Equal(t, 4, len(response.Tokens))

	assert.Equal(t, models.Token{
		RawText:   "hi ",
		TokenType: models.TextTokenType,
	}, response.Tokens[0])

	assert.Equal(t, models.Token{
		RawText:      "$(user)",
		TokenType:    models.VariableTokenType,
		VariableName: "user",
		Arguments: []models.Token{
			models.Token{
				RawText:   "user",
				TokenType: models.TextTokenType,
			},
		},
	}, response.Tokens[1])

	assert.Equal(t, models.Token{
		RawText:   " hello ",
		TokenType: models.TextTokenType,
	}, response.Tokens[2])

	assert.Equal(t, models.Token{
		RawText:      "$(user)",
		TokenType:    models.VariableTokenType,
		VariableName: "user",
		Arguments: []models.Token{
			models.Token{
				RawText:   "user",
				TokenType: models.TextTokenType,
			},
		},
	}, response.Tokens[3])
}

// response with multiple variables of same type
// "Display name of $(userid) is $(display_name)"
func TestVariableMultipleVariables(t *testing.T) {
	response := ParseResponse("hi $(user) to $(channel)")

	assert.Equal(t, "hi $(user) to $(channel)", response.RawText)
	assert.Equal(t, 4, len(response.Tokens))

	assert.Equal(t, models.Token{
		RawText:   "hi ",
		TokenType: models.TextTokenType,
	}, response.Tokens[0])

	assert.Equal(t, models.Token{
		RawText:      "$(user)",
		TokenType:    models.VariableTokenType,
		VariableName: "user",
		Arguments: []models.Token{
			models.Token{
				RawText:   "user",
				TokenType: models.TextTokenType,
			},
		},
	}, response.Tokens[1])

	assert.Equal(t, models.Token{
		RawText:   " to ",
		TokenType: models.TextTokenType,
	}, response.Tokens[2])

	assert.Equal(t, models.Token{
		RawText:      "$(channel)",
		TokenType:    models.VariableTokenType,
		VariableName: "channel",
		Arguments: []models.Token{
			models.Token{
				RawText:   "channel",
				TokenType: models.TextTokenType,
			},
		},
	}, response.Tokens[3])
}

func TestMultipleVariablesNoSpace(t *testing.T) {
	response := ParseResponse("A$(user)B$(channel)C")

	assert.Equal(t, "A$(user)B$(channel)C", response.RawText)
	assert.Equal(t, 5, len(response.Tokens))

	assert.Equal(t, models.Token{
		RawText:   "A",
		TokenType: models.TextTokenType,
	}, response.Tokens[0])

	assert.Equal(t, models.Token{
		RawText:      "$(user)",
		TokenType:    models.VariableTokenType,
		VariableName: "user",
		Arguments: []models.Token{
			models.Token{
				RawText:   "user",
				TokenType: models.TextTokenType,
			},
		},
	}, response.Tokens[1])

	assert.Equal(t, models.Token{
		RawText:   "B",
		TokenType: models.TextTokenType,
	}, response.Tokens[2])

	assert.Equal(t, models.Token{
		RawText:      "$(channel)",
		TokenType:    models.VariableTokenType,
		VariableName: "channel",
		Arguments: []models.Token{
			models.Token{
				RawText:   "channel",
				TokenType: models.TextTokenType,
			},
		},
	}, response.Tokens[3])

	assert.Equal(t, models.Token{
		RawText:   "C",
		TokenType: models.TextTokenType,
	}, response.Tokens[4])
}

// response with multiple variables of different types
// "$(user) has followed for $(follow_age $(user))"
func TestNestedVariable(t *testing.T) {
	response := ParseResponse("hi $(urlfetch http://twitch.tv/$(user)/01)")

	assert.Equal(t, "hi $(urlfetch http://twitch.tv/$(user)/01)", response.RawText)
	assert.Equal(t, 2, len(response.Tokens))

	assert.Equal(t, models.Token{
		RawText:   "hi ",
		TokenType: models.TextTokenType,
	}, response.Tokens[0])

	vToken := response.Tokens[1]
	assert.Equal(t, "$(urlfetch http://twitch.tv/$(user)/01)", vToken.RawText)
	assert.Equal(t, models.VariableTokenType, vToken.TokenType)
	assert.Equal(t, "urlfetch", vToken.VariableName)
	assert.Equal(t, 3, len(vToken.Arguments))

	// Examine nested tokens
	assert.Equal(t, models.Token{
		RawText:   "urlfetch http://twitch.tv/",
		TokenType: models.TextTokenType,
	}, vToken.Arguments[0])
	assert.Equal(t, models.Token{
		RawText:      "$(user)",
		TokenType:    models.VariableTokenType,
		VariableName: "user",
		Arguments: []models.Token{
			models.Token{
				RawText:   "user",
				TokenType: models.TextTokenType,
			},
		},
	}, vToken.Arguments[1])
	assert.Equal(t, models.Token{
		RawText:   "/01",
		TokenType: models.TextTokenType,
	}, vToken.Arguments[2])
}

// More complicated nested variables
func TestNestedVariable2(t *testing.T) {
	response := ParseResponse("$(user) has followed for $(follow_age $(user))")

	assert.Equal(t, "$(user) has followed for $(follow_age $(user))", response.RawText)
	assert.Equal(t, 3, len(response.Tokens))

	assert.Equal(t, models.Token{
		RawText:      "$(user)",
		TokenType:    models.VariableTokenType,
		VariableName: "user",
		Arguments: []models.Token{
			models.Token{
				RawText:   "user",
				TokenType: models.TextTokenType,
			},
		},
	}, response.Tokens[0])
	assert.Equal(t, models.Token{
		RawText:   " has followed for ",
		TokenType: models.TextTokenType,
	}, response.Tokens[1])

	vToken := response.Tokens[2]
	assert.Equal(t, "$(follow_age $(user))", vToken.RawText)
	assert.Equal(t, models.VariableTokenType, vToken.TokenType)
	assert.Equal(t, "follow_age", vToken.VariableName)
	assert.Equal(t, 2, len(vToken.Arguments))

	// Examine nested tokens
	assert.Equal(t, models.Token{
		RawText:   "follow_age ",
		TokenType: models.TextTokenType,
	}, vToken.Arguments[0])
	assert.Equal(t, models.Token{
		RawText:      "$(user)",
		TokenType:    models.VariableTokenType,
		VariableName: "user",
		Arguments: []models.Token{
			models.Token{
				RawText:   "user",
				TokenType: models.TextTokenType,
			},
		},
	}, vToken.Arguments[1])
}

func TestContinuousNestedVariables(t *testing.T) {
	response := ParseResponse("$(display_name)$(user)$(follow_age $(user))")

	assert.Equal(t, "$(display_name)$(user)$(follow_age $(user))", response.RawText)
	assert.Equal(t, 3, len(response.Tokens))

	assert.Equal(t, models.Token{
		RawText:      "$(display_name)",
		TokenType:    models.VariableTokenType,
		VariableName: "display_name",
		Arguments: []models.Token{
			models.Token{
				RawText:   "display_name",
				TokenType: models.TextTokenType,
			},
		},
	}, response.Tokens[0])
	assert.Equal(t, models.Token{
		RawText:      "$(user)",
		TokenType:    models.VariableTokenType,
		VariableName: "user",
		Arguments: []models.Token{
			models.Token{
				RawText:   "user",
				TokenType: models.TextTokenType,
			},
		},
	}, response.Tokens[1])

	vToken := response.Tokens[2]
	assert.Equal(t, "$(follow_age $(user))", vToken.RawText)
	assert.Equal(t, models.VariableTokenType, vToken.TokenType)
	assert.Equal(t, "follow_age", vToken.VariableName)
	assert.Equal(t, 2, len(vToken.Arguments))

	// Examine nested tokens
	assert.Equal(t, models.Token{
		RawText:   "follow_age ",
		TokenType: models.TextTokenType,
	}, vToken.Arguments[0])
	assert.Equal(t, models.Token{
		RawText:      "$(user)",
		TokenType:    models.VariableTokenType,
		VariableName: "user",
		Arguments: []models.Token{
			models.Token{
				RawText:   "user",
				TokenType: models.TextTokenType,
			},
		},
	}, vToken.Arguments[1])
}

func TestDeeplyNestedVariables1(t *testing.T) {
	response := ParseResponse("$($(a $(b) c))")

	assert.Equal(t, "$($(a $(b) c))", response.RawText)
	assert.Equal(t, 1, len(response.Tokens))

	vToken := response.Tokens[0]
	assert.Equal(t, "$($(a $(b) c))", vToken.RawText)
	assert.Equal(t, models.VariableTokenType, vToken.TokenType)
	assert.Equal(t, "", vToken.VariableName) // Note that this is invalid token, therefore empty variable name.

	nestedTokens1 := vToken.Arguments
	assert.Equal(t, 1, len(nestedTokens1))
	assert.Equal(t, "$(a $(b) c)", nestedTokens1[0].RawText)
	assert.Equal(t, models.VariableTokenType, nestedTokens1[0].TokenType)
	assert.Equal(t, "a", nestedTokens1[0].VariableName)

	nestedTokens2 := nestedTokens1[0].Arguments
	assert.Equal(t, 3, len(nestedTokens2))

	// Examine nested tokens
	assert.Equal(t, models.Token{
		RawText:   "a ",
		TokenType: models.TextTokenType,
	}, nestedTokens2[0])
	assert.Equal(t, models.Token{
		RawText:      "$(b)",
		TokenType:    models.VariableTokenType,
		VariableName: "b",
		Arguments: []models.Token{
			models.Token{
				RawText:   "b",
				TokenType: models.TextTokenType,
			},
		},
	}, nestedTokens2[1])
	assert.Equal(t, models.Token{
		RawText:   " c",
		TokenType: models.TextTokenType,
	}, nestedTokens2[2])
}

// When starting and ending tags for variables do not match
// Test when a variable is unfinished.
// In such case, validation will fail because of malformed response text.
func TestUnfinishedNestedVariable1(t *testing.T) {
	response := ParseResponse("followed for $(follow_age $(user) . Thanks")

	assert.Equal(t, "followed for $(follow_age $(user) . Thanks", response.RawText)
	// " . Thanks" belongs to $(follow_age) variable because of malformed response.
	assert.Equal(t, 2, len(response.Tokens))

	assert.Equal(t, models.Token{
		RawText:   "followed for ",
		TokenType: models.TextTokenType,
	}, response.Tokens[0])

	// The second token is complicated
	vToken := response.Tokens[1]
	assert.Equal(t, "$(follow_age $(user) . Thanks", vToken.RawText)
	assert.Equal(t, models.VariableTokenType, vToken.TokenType)
	assert.Equal(t, "follow_age", vToken.VariableName)

	// Inspect nested tokens
	nestedTokens1 := vToken.Arguments
	assert.Equal(t, 3, len(nestedTokens1))
	assert.Equal(t, models.Token{
		RawText:   "follow_age ",
		TokenType: models.TextTokenType,
	}, nestedTokens1[0])
	assert.Equal(t, models.Token{
		RawText:      "$(user)",
		TokenType:    models.VariableTokenType,
		VariableName: "user",
		Arguments: []models.Token{
			models.Token{
				RawText:   "user",
				TokenType: models.TextTokenType,
			},
		},
	}, nestedTokens1[1])
	assert.Equal(t, models.Token{
		RawText:   " . Thanks",
		TokenType: models.TextTokenType,
	}, nestedTokens1[2])
}

func TestUnfinishedNestedVariable2(t *testing.T) {
	response := ParseResponse("followed for $(follow_age $(user . Thanks)")

	assert.Equal(t, "followed for $(follow_age $(user . Thanks)", response.RawText)
	// from $(follow_age till the end of string is the second token
	assert.Equal(t, 2, len(response.Tokens))

	assert.Equal(t, models.Token{
		RawText:   "followed for ",
		TokenType: models.TextTokenType,
	}, response.Tokens[0])

	// The second token is complicated
	vToken := response.Tokens[1]
	assert.Equal(t, "$(follow_age $(user . Thanks)", vToken.RawText)
	assert.Equal(t, models.VariableTokenType, vToken.TokenType)
	assert.Equal(t, "follow_age", vToken.VariableName)

	// Inspect nested tokens
	nestedTokens1 := vToken.Arguments
	assert.Equal(t, 2, len(nestedTokens1))
	assert.Equal(t, models.Token{
		RawText:   "follow_age ",
		TokenType: models.TextTokenType,
	}, nestedTokens1[0])
	assert.Equal(t, models.Token{
		RawText:      "$(user . Thanks)",
		TokenType:    models.VariableTokenType,
		VariableName: "user",
		Arguments: []models.Token{
			models.Token{
				RawText:   "user . Thanks",
				TokenType: models.TextTokenType,
			},
		},
	}, nestedTokens1[1])
}

// When the opening tag for variable is malformed, then it is a text response.
func TestVariableWithPartialOpeningTag(t *testing.T) {
	response := ParseResponse("hello $user)")

	assert.Equal(t, "hello $user)", response.RawText)
	assert.Equal(t, 1, len(response.Tokens))
	assert.Equal(t, models.Token{
		RawText:   "hello $user)",
		TokenType: models.TextTokenType,
	}, response.Tokens[0])
}

func TestFullVariableWithNoName(t *testing.T) {
	response := ParseResponse("hello $() world")

	assert.Equal(t, "hello $() world", response.RawText)
	assert.Equal(t, 3, len(response.Tokens))

	tokens := response.Tokens

	assert.Equal(t, models.Token{
		RawText:   "hello ",
		TokenType: models.TextTokenType,
	}, tokens[0])
	// Note that in this case, no arguments exist
	assert.Equal(t, models.Token{
		RawText:      "$()",
		TokenType:    models.VariableTokenType,
		VariableName: "",
	}, tokens[1])
	assert.Equal(t, models.Token{
		RawText:   " world",
		TokenType: models.TextTokenType,
	}, tokens[2])
}

func TestPartialVariableWithNoName(t *testing.T) {
	response := ParseResponse("$(")

	assert.Equal(t, "$(", response.RawText)
	assert.Equal(t, 1, len(response.Tokens))

	// Note that in this case, no arguments exist
	assert.Equal(t, models.Token{
		RawText:   "$(",
		TokenType: models.VariableTokenType,
	}, response.Tokens[0])
}

// ---------------------------------------------------------------------------------
// ---------------------------------------------------------------------------------
// TODO: Create tokens, instead of parsed response, as input for GetVariableName()
// ---------------------------------------------------------------------------------
// ---------------------------------------------------------------------------------

func TestGetVariableNameNormal(t *testing.T) {
	response := ParseResponse("$(user)")
	assert.Equal(t, "user", response.Tokens[0].VariableName)
}

func TestGetVariableNameWithNoName(t *testing.T) {
	response := ParseResponse("  $()  ")
	assert.Equal(t, "", response.Tokens[1].VariableName)

	response = ParseResponse("  $(    )  ")
	assert.Equal(t, "", response.Tokens[1].VariableName)

	// Check non-space whitespace characters
	response = ParseResponse("  $(\t  )  ")
	assert.Equal(t, "", response.Tokens[1].VariableName)
}

// Test GetVariableName() with different whitespaces
func TestGetVariableNameWithWhitespaces(t *testing.T) {
	response := ParseResponse("$(  user   )")
	assert.Equal(t, "user", response.Tokens[0].VariableName)

	// Note the tab character '\t' between "user" and "ee"
	response = ParseResponse("   d $( \t user\tee )  ")
	assert.Equal(t, "user", response.Tokens[1].VariableName)
}

// Variable error cases..

// Test isStartOfVariable(), in various cases
