package commandplugins

import (
	"testing"

	models "github.com/c-rainbow/simplechatbot/models"
	plugins "github.com/c-rainbow/simplechatbot/plugins"
	chatplugins "github.com/c-rainbow/simplechatbot/plugins/chat"
	mock_repository "github.com/c-rainbow/simplechatbot/repository/mock"
	twitch_irc "github.com/gempir/go-twitch-irc"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

const (
	testChannel     = "test_channel"
	testPluginType  = "TestPluginType"
	testCommandName = "!test"
)

var (
	testSender  = twitch_irc.User{DisplayName: "TestUser"}
	testMessage = twitch_irc.Message{Text: "Hello this is test message"}
	testCommand = models.Command{
		Name:       testCommandName,
		Enabled:    true,
		PluginType: testPluginType,
		Permission: chatplugins.PermissionEveryone,
	}
)

func TestValidateInputsNilCommand(t *testing.T) {
	err := ValidateInputs(nil, testChannel, testPluginType, &testSender, &testMessage)
	assert.Equal(t, CommandNotFoundError, err)
}

// Test ValidateInputs() when plugin type does not match
func TestValidateInputsIncorrectPluginType(t *testing.T) {
	err := ValidateInputs(&testCommand, testChannel, "InvalidPluginType", &testSender, &testMessage)
	assert.Equal(t, plugins.ErrIncorrectPluginType, err)
}

// Test ValidateInputs() when command with the name is not found
func TestValidateInputsNoPermission(t *testing.T) {
	command := testCommand
	command.Permission = chatplugins.PermissionModerator
	err := ValidateInputs(&command, testChannel, testPluginType, &testSender, &testMessage)
	assert.Equal(t, NoPermissionError, err)
}

func Testddddd(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockRepo := mock_repository.NewMockSingleBotRepositoryT(ctrl)
	mockRepo.EXPECT().GetCommandByChannelAndName(gomock.Eq(testChannel), gomock.Eq(testCommandName)).
		Return(nil).AnyTimes()
}
