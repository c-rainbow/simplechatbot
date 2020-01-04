package common

import (
	"testing"

	models "github.com/c-rainbow/simplechatbot/models"
	parser "github.com/c-rainbow/simplechatbot/parser"
	plugins "github.com/c-rainbow/simplechatbot/plugins"
	chatplugins "github.com/c-rainbow/simplechatbot/plugins/chat"
	twitch_irc "github.com/gempir/go-twitch-irc"
	"github.com/stretchr/testify/assert"
)

const (
	testChannel1    = "testchannel"
	testChannel2    = "testchannel2"
	testPluginType  = "TestPluginType"
	testCommandName = "!test"
)

var (
	testSender  = twitch_irc.User{Name: testChannel1, DisplayName: "TestChannel1"}
	testMessage = twitch_irc.PrivateMessage{
		Message: "Hello this is test message",
		Tags:    map[string]string{},
	}
	testArgs    = []string{"testArg"}
	testCommand = models.Command{
		Name:       testCommandName,
		Enabled:    true,
		PluginType: testPluginType,
		Permission: models.PermissionEveryone,
		Responses: map[string]models.ParsedResponse{
			models.DefaultResponseKey: testResponse,
		},
	}
	testResponse = models.ParsedResponse{
		RawText: "$(arg1)",
		Tokens: []models.Token{
			{
				RawText:      "$(arg1)",
				TokenType:    models.VariableTokenType,
				VariableName: "arg1",
				Arguments: []models.Token{
					{
						RawText:   "arg1",
						TokenType: models.TextTokenType,
					},
				},
			},
		},
	}
)

func TestValidateBasicInputsNilCommand(t *testing.T) {
	err := ValidateBasicInputs(nil, testChannel1, testPluginType, &testSender, &testMessage)
	assert.Equal(t, chatplugins.ErrCommandNotFound, err)
}

// Test ValidateInputs() when plugin type does not match
func TestValidateBasicInputsIncorrectPluginType(t *testing.T) {
	err := ValidateBasicInputs(&testCommand, testChannel1, "InvalidPluginType", &testSender, &testMessage)
	assert.Equal(t, plugins.ErrIncorrectPluginType, err)
}

// Test ValidateInputs() when command with the name is not found
func TestValidateBasicInputsNoPermission(t *testing.T) {
	command := testCommand
	command.Permission = models.PermissionModerator
	err := ValidateBasicInputs(&command, "another_channel", testPluginType, &testSender, &testMessage)
	assert.Equal(t, chatplugins.ErrNoPermission, err)
}

// Test ConvertToResponseText() without matching response for the response key
func TestConvertToResponseTextNoResponse(t *testing.T) {
	message, err := ConvertToResponseText(
		&testCommand, "InvalidResponseKey", testChannel1, &testSender, &testMessage, testArgs)
	assert.Empty(t, message)
	assert.Equal(t, chatplugins.ErrResponseNotFound, err)
}

// Test ConvertToResponseText() without matching response for the response key
func TestConvertToResponseTextConvertError(t *testing.T) {
	// The test response has $(arg1) but testArgs has only 1 element.
	message, err := ConvertToResponseText(
		&testCommand, models.DefaultResponseKey, testChannel1, &testSender, &testMessage, testArgs)
	assert.Empty(t, message)
	assert.Equal(t, parser.InvalidVariableNameError, err)
}

func TestConvertToResponseSuccess(t *testing.T) {
	args := []string{"arg0", "This is arg1"}
	message, err := ConvertToResponseText(
		&testCommand, models.DefaultResponseKey, testChannel1, &testSender, &testMessage, args)
	assert.Nil(t, err)
	assert.Equal(t, "This is arg1", message)

}

func TestUserPermissionsStreamer(t *testing.T) {
	// Streamer always has permission, regardless of command.
	assert.True(t, UserHasPermission(testSender.Name, nil, &testSender, nil))
}

func TestUserPermissionsEveryone(t *testing.T) {
	command := testCommand // Copy testCommand not to make changes to the original one.

	// Everyone has permission if permission bit is set, regardless of channel and message
	command.Permission = models.PermissionEveryone
	assert.True(t, UserHasPermission(testChannel2, &command, &testSender, nil))

	// Same whenever other bits are set.
	command.Permission = models.PermissionEveryone | models.PermissionSubscriber
	assert.True(t, UserHasPermission(testChannel2, &command, &testSender, nil))
}

func TestUserPermissionsSubscriber(t *testing.T) {
	message := testMessage // Copy not to make changes to the original one.
	command := testCommand
	command.Permission = models.PermissionModerator | models.PermissionSubscriber

	// When subscriber tag value is empty
	message.Tags["subscriber"] = ""
	assert.False(t, UserHasPermission(testChannel2, &command, &testSender, &message))

	// When subscriber tag value is "1"
	message.Tags["subscriber"] = "1"
	assert.True(t, UserHasPermission(testChannel2, &command, &testSender, &message))

	// When subscriber tag value is "0"
	message.Tags["subscriber"] = "0"
	assert.False(t, UserHasPermission(testChannel2, &command, &testSender, &message))
}

func TestUserPermissionsModerator(t *testing.T) {
	command := testCommand // Copy testCommand not to make changes to the original one.
	command.Permission = models.PermissionModerator

	message := testMessage // Copy testMessage not to make changes to the original one.

	// When moderator tag value is empty
	message.Tags["mod"] = ""
	assert.False(t, UserHasPermission(testChannel2, &command, &testSender, &message))

	// When moderator tag value is "1"
	message.Tags["mod"] = "1"
	assert.True(t, UserHasPermission(testChannel2, &command, &testSender, &message))

	// When moderator tag value is "0"
	message.Tags["mod"] = "0"
	assert.False(t, UserHasPermission(testChannel2, &command, &testSender, &message))

	// Make sure than non-mod subscriber has no permission
	message.Tags["subscriber"] = "1"
	assert.False(t, UserHasPermission(testChannel2, &command, &testSender, &message))

}
