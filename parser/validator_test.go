package parser

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetVariableName(t *testing.T) {
	response := ParseResponse("$(user)")
	//response.
	_ = response
}

func TestValidateEmpty(t *testing.T) {
	response := ParseResponse("")
	err := Validate(response)
	assert.Equal(t, EmptyResponseError, err)
}

func TestValidateSimpleText(t *testing.T) {
	response := ParseResponse("hello world")
	err := Validate(response)
	assert.Nil(t, err)
}

func TestValidateEntireVariable(t *testing.T) {
	response := ParseResponse("$(user)")
	err := Validate(response)
	assert.Nil(t, err)
}

func TestValidateSpaceAroundVariable(t *testing.T) {
	response := ParseResponse("  $(user)  ")
	err := Validate(response)
	assert.Nil(t, err)
}

func TestValidateSpaceAroundVariableName(t *testing.T) {
	response := ParseResponse("  $(  user   )  ")
	err := Validate(response)
	assert.Nil(t, err)
}

func TestValidateNoVariableName(t *testing.T) {
	response := ParseResponse("  $()  ")
	err := Validate(response)
	assert.Equal(t, VariableNameNotFoundError, err)

	response = ParseResponse("  $(    )  ")
	err = Validate(response)
	assert.Equal(t, VariableNameNotFoundError, err)

	// Check non-space whitespace characters
	response = ParseResponse("  $(\t  )  ")
	err = Validate(response)
	assert.Equal(t, VariableNameNotFoundError, err)
}

func TestValidateNoClosingTag(t *testing.T) {
	response := ParseResponse("  $(user ")
	err := Validate(response)
	assert.Equal(t, VariableHasNoCloseTagError, err)
}

func TestValidateMultipleVariables(t *testing.T) {
	response := ParseResponse("hi $(user) hello $(user)")
	err := Validate(response)
	assert.Nil(t, err)

	// Different variables
	response = ParseResponse("A$(user)B$(channel)C")
	err = Validate(response)
	assert.Nil(t, err)
}

func TestValidateCorrectNestedVariables(t *testing.T) {
	response := ParseResponse("$(user) has followed for $(follow_age $(user))")
	err := Validate(response)
	assert.Nil(t, err)
}

func TestValidateVariablesWithNoName(t *testing.T) {
	response := ParseResponse("  $()  ")
	err := Validate(response)
	assert.Equal(t, VariableNameNotFoundError, err)

	response = ParseResponse("  $(    )  ")
	err = Validate(response)
	assert.Equal(t, VariableNameNotFoundError, err)

	// Check non-space whitespace characters
	response = ParseResponse("  $(\t  )  ")
	err = Validate(response)
	assert.Equal(t, VariableNameNotFoundError, err)
}

// Test GetVariableName() with different whitespaces
func TestValidateInvalidName(t *testing.T) {
	response := ParseResponse("$(invalid)")
	err := Validate(response)
	assert.Equal(t, VariableNameNotAllowedError, err)

	response = ParseResponse("Nested invalid $(follow_age $(invalid)) variable check ")
	err = Validate(response)
	assert.Equal(t, VariableNameNotAllowedError, err)
}

func TestValidateIncorrectNestedVariables(t *testing.T) {
	response := ParseResponse("$(user) has followed for $(follow_age $(uptime))")
	err := Validate(response)
	assert.Equal(t, VariableCanNotBeNestedError, err)
}
