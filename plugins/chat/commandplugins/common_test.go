package commandplugins

import (
	"testing"

	models "github.com/c-rainbow/simplechatbot/models"
	plugins "github.com/c-rainbow/simplechatbot/plugins"
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
	testSender  = &twitch_irc.User{DisplayName: "TestUser"}
	testMessage = &twitch_irc.Message{Text: "Hello this is test message"}
	testCommand = &models.Command{
		Name:    testCommandName,
		Enabled: true,
	}
)

// Test CommonRead() when command with the name is not found
func TestCommonReadCommandNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockRepo := mock_repository.NewMockSingleBotRepositoryT(ctrl)
	mockRepo.EXPECT().GetCommandByChannelAndName(gomock.Eq(testChannel), gomock.Eq(testCommandName)).
		Return(nil).AnyTimes()
	command, err := CommonRead(
		mockRepo, testCommandName, testChannel, testPluginType, testSender, testMessage)

	assert.Nil(t, command)
	assert.Equal(t, CommandNotFoundError, err)
}

// Test CommonRead() when plugin type does not match
func TestCommonReadIncorrectPluginType(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockRepo := mock_repository.NewMockSingleBotRepositoryT(ctrl)
	mockRepo.EXPECT().GetCommandByChannelAndName(gomock.Eq(testChannel), gomock.Eq(testCommandName)).
		Return(testCommand).AnyTimes()
	command, err := CommonRead(
		mockRepo, testCommandName, testChannel, testPluginType, testSender, testMessage)

	assert.Nil(t, command)
	assert.Equal(t, plugins.IncorrectPluginTypeError, err)
}
