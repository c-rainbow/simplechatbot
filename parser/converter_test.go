package parser

import (
	"fmt"
	"testing"

	twitch_irc "github.com/gempir/go-twitch-irc"
	"github.com/stretchr/testify/assert"
)

func TestConvertNoVariable(t *testing.T) {
	response := ParseResponse("hello")
	message, err := ConvertResponse(response, "test_channel", nil, nil)
	assert.Nil(t, err)
	assert.Equal(t, "hello", message)
}

func TestConvertResponseSimple(t *testing.T) {
	response := ParseResponse("hello $(user), welcome to $(channel) ")

	// Create sample user for testing
	user := twitch_irc.User{Username: "1234", DisplayName: "TestUser"}
	message, err := ConvertResponse(response, "test_channel", &user, nil)
	assert.Nil(t, err)
	assert.Equal(t, "hello TestUser, welcome to test_channel ", message)
}

func TestConvertVariableNotEnabled1(t *testing.T) {
	response := ParseResponse("The streamer is playing $(game)")
	message, err := ConvertResponse(response, "test_channel", nil, nil)
	fmt.Println("message: ", message)
	assert.Equal(t, VariableNotEnabledError, err)
}

func TestConvertVariableNotEnabled2(t *testing.T) {
	response := ParseResponse("$(user) has followed for $(follow_age $(user))")
	user := twitch_irc.User{Username: "1234", DisplayName: "TestUser"}
	_, err := ConvertResponse(response, "test_channel", &user, nil)
	assert.Equal(t, VariableNotEnabledError, err)
}
