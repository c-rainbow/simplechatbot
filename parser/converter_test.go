package parser

import (
	"testing"

	twitch_irc "github.com/gempir/go-twitch-irc"
	"github.com/stretchr/testify/assert"
)

func TestConvertNoVariable(t *testing.T) {
	response := ParseResponse("hello")
	message, err := ConvertResponse(response, "test_channel", nil, nil, nil, nil)
	assert.Nil(t, err)
	assert.Equal(t, "hello", message)
}

func TestConvertResponseSimple(t *testing.T) {
	response := ParseResponse("hello $(user), welcome to $(channel) ")

	// Create sample user for testing
	user := twitch_irc.User{Name: "1234", DisplayName: "TestUser"}
	message, err := ConvertResponse(response, "test_channel", &user, nil, nil, nil)
	assert.Nil(t, err)
	assert.Equal(t, "hello TestUser, welcome to test_channel ", message)
}

/*
func TestConvertVariableNotEnabled1(t *testing.T) {
	response := ParseResponse("The streamer is playing $(game)")
	_, err := ConvertResponse(response, "test_channel", nil, nil, nil, nil)
	assert.Equal(t, VariableNotEnabledError, err)
}

func TestConvertVariableNotEnabled2(t *testing.T) {
	response := ParseResponse("$(user) has followed for $(follow_age $(user))")
	user := twitch_irc.User{Name: "1234", DisplayName: "TestUser"}
	_, err := ConvertResponse(response, "test_channel", &user, nil, nil, nil)
	assert.Equal(t, VariableNotEnabledError, err)
}
*/

func TestConvertArguments(t *testing.T) {
	response := ParseResponse("hi $(arg0) $(arg1) $(arg2) $(arg3) $(arg4) $(arg5) $(arg0) bye")
	args := []string{"Zero", "one", "two", "three", "four", "five"}
	message, err := ConvertResponse(response, "test_channel", nil, nil, nil, args)
	assert.Nil(t, err)
	assert.Equal(t, "hi Zero one two three four five Zero bye", message)
}

func TestConvertVariablesMultipleTypes(t *testing.T) {
	response := ParseResponse("hi $(user) $(arg0) bye")
	user := twitch_irc.User{Name: "1234", DisplayName: "TestUser"}
	args := []string{"This is argument zero"}
	message, err := ConvertResponse(response, "test_channel", &user, nil, nil, args)
	assert.Nil(t, err)
	assert.Equal(t, "hi TestUser This is argument zero bye", message)
}

func TestConvertUnsupportedType(t *testing.T) {
	response := ParseResponse("$(invalid)")
	_, err := ConvertResponse(response, "", nil, nil, nil, nil)
	assert.Equal(t, UnsupportedVariableTypeError, err)
}
