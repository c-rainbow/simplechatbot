package parser

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPass(t *testing.T) {

}

func TestGetVariableName(t *testing.T) {
	// TODO: This is a bad habit to depend on ParseResponse in another module.
	response := ParseResponse("$(user)")
	//response.
	_ = response
}

func TestValidateEmpty(t *testing.T) {
	// TODO: This is a bad habit to depend on ParseResponse in another module.
	response := ParseResponse("")
	err := Validate(response)
	assert.Equal(t, EmptyResponseError, err)
}

func TestValidateSimpleText(t *testing.T) {
	// TODO: This is a bad habit to depend on ParseResponse in another module.
	response := ParseResponse("hello world")
	err := Validate(response)
	assert.Nil(t, err)
}

func TestValidateEntireVariable(t *testing.T) {
	// TODO: This is a bad habit to depend on ParseResponse in another module.
	response := ParseResponse("$(user)")
	err := Validate(response)
	assert.Nil(t, err)
}

func TestValidateSpaceAroundVariable(t *testing.T) {
	// TODO: This is a bad habit to depend on ParseResponse in another module.
	response := ParseResponse("  $(user)  ")
	err := Validate(response)
	assert.Nil(t, err)
}

func TestValidateSpaceAroundVariableName(t *testing.T) {
	// TODO: This is a bad habit to depend on ParseResponse in another module.
	response := ParseResponse("  $(  user   )  ")
	err := Validate(response)
	assert.Nil(t, err)
}

func TestValidateNoVariableName(t *testing.T) {
	// TODO: This is a bad habit to depend on ParseResponse in another module.
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
	// TODO: This is a bad habit to depend on ParseResponse in another module.
	response := ParseResponse("hi $(user) hello $(user)")
	err := Validate(response)
	assert.Nil(t, err)

	// Different variables
	response = ParseResponse("A$(user)B$(channel)C")
	err = Validate(response)
	assert.Nil(t, err)
}

func TestValidateCorrectNestedVariables(t *testing.T) {
	// TODO: This is a bad habit to depend on ParseResponse in another module.
	response := ParseResponse("$(user) has followed for $(follow_age $(user))")
	err := Validate(response)
	assert.Nil(t, err)

}

// TODO: Add more tests for GetVariableName

func TestGetVariableNameWithArguments(t *testing.T) {
	response := ParseResponse("$(follow_age $(user))")
	vName, err := GetVariableName(&response.Tokens[0])
	assert.Nil(t, err)
	assert.Equal(t, "follow_age", vName)
}

func TestGetVariableNameNormal(t *testing.T) {
	response := ParseResponse("$(user)")
	vName, err := GetVariableName(&response.Tokens[0])
	assert.Nil(t, err)
	assert.Equal(t, "user", vName)
}

func TestGetVariableNameWithNoName(t *testing.T) {
	// TODO: This is a bad habit to depend on ParseResponse in another module.
	response := ParseResponse("  $()  ")
	_, err := GetVariableName(&response.Tokens[1])
	assert.Equal(t, VariableNameNotFoundError, err)

	response = ParseResponse("  $(    )  ")
	_, err = GetVariableName(&response.Tokens[1])
	assert.Equal(t, VariableNameNotFoundError, err)

	// Check non-space whitespace characters
	response = ParseResponse("  $(\t  )  ")
	_, err = GetVariableName(&response.Tokens[1])
	assert.Equal(t, VariableNameNotFoundError, err)
}

// Test GetVariableName() with different whitespaces
func TestGetVariableNameWithWhitespaces(t *testing.T) {
	response := ParseResponse("$(  user   )")
	vName, err := GetVariableName(&response.Tokens[0])
	assert.Nil(t, err)
	assert.Equal(t, "user", vName)

	// Note the tab character '\t' between "user" and "ee"
	response = ParseResponse("   d $( \t user\tee )  ")
	vName, err = GetVariableName(&response.Tokens[1])
	assert.Nil(t, err)
	assert.Equal(t, "user", vName)
}

// TODO: Enable this test after nestability check
/*func TestValidateIncorrectNestedVariables(t *testing.T) {
	response := ParseResponse("$(user) has followed for $(follow_age $(uptime))")
	err := Validate(response)
	assert.Equal(t, VariableCanNotBeNestedError, err)
}*/
