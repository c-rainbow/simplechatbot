package commandplugins

import (
	"github.com/c-rainbow/simplechatbot/client"
	models "github.com/c-rainbow/simplechatbot/models"
	chatplugins "github.com/c-rainbow/simplechatbot/plugins/chat"
	"github.com/c-rainbow/simplechatbot/plugins/chat/common"
	"github.com/c-rainbow/simplechatbot/repository"
	twitch_irc "github.com/gempir/go-twitch-irc"
)

const (
	// DeleteCommandPluginType plugin type name to delete an existing command (of any type)
	DeleteCommandPluginType = "DeleteCommandPluginType"
)

// DeleteCommandPlugin to delete an existing command (of any type)
type DeleteCommandPlugin struct {
	ircClient client.TwitchClientT
	repo      repository.SingleBotRepositoryT
}

var _ chatplugins.ChatCommandPluginT = (*DeleteCommandPlugin)(nil)

// NewDeleteCommandPlugin creates a new DeleteCommandPlugin
func NewDeleteCommandPlugin(
	ircClient client.TwitchClientT, repo repository.SingleBotRepositoryT) chatplugins.ChatCommandPluginT {
	return &DeleteCommandPlugin{ircClient: ircClient, repo: repo}
}

// GetPluginType returns plugin type
func (plugin *DeleteCommandPlugin) GetPluginType() string {
	return DeleteCommandPluginType
}

// ReactToChat reacts to chat
func (plugin *DeleteCommandPlugin) ReactToChat(
	command *models.Command, channel string, sender *twitch_irc.User, message *twitch_irc.PrivateMessage) {
	var err error
	var targetCommand *models.Command

	targetCommandName, _ := GetTargetCommandNameAndResponse(message.Message)

	// TODO: Is it possible to get away from this continuous err == nil check?
	err = common.ValidateBasicInputs(command, channel, DeleteCommandPluginType, sender, message)
	if err == nil {
		targetCommand, err = getTargetCommand(channel, targetCommandName, plugin.repo)
	}
	if err == nil {
		err = plugin.validateTargetCommand(targetCommand)
	}
	if err == nil {
		err = plugin.repo.DeleteCommand(channel, targetCommand)
	}

	responseText, err := plugin.GetResponseText(command, targetCommand, channel, sender, message, err)

	common.SendToChatClient(plugin.ircClient, channel, responseText)
	common.HandleError(err)
}

// This function is slightly different between add/edit/delete command. Hard to merge into a common function.
func (plugin *DeleteCommandPlugin) validateTargetCommand(targetCommand *models.Command) error {
	// Can't delete non-existing command
	if targetCommand == nil {
		return chatplugins.ErrTargetCommandNotFound
	}
	return nil
}

// GetResponseText gets response text of the executed command, based on the errors and progress so far.
func (plugin *DeleteCommandPlugin) GetResponseText(
	command *models.Command, targetCommand *models.Command, channel string, sender *twitch_irc.User,
	message *twitch_irc.PrivateMessage, err error) (string, error) {
	// Get response key, build args, get parsed response, and convert it to text
	responseKey := plugin.GetResponseKey(err)
	args := []string{targetCommand.Name}
	return common.ConvertToResponseText(command, responseKey, channel, sender, message, args)
}

// GetResponseKey returns response key from error type to build response text accordingly.
func (plugin *DeleteCommandPlugin) GetResponseKey(err error) string {
	// Normal case.
	if err == nil {
		return models.DefaultResponseKey
	}

	/*
		Failure case.
		Design decision: We can return different messages per error type in two ways
		(1) switch statement with each known error cases, manually assigning response key, like the code below
		(2) Each error has unique error type string, and we use it as response key. For example,
			parsedResponse, exists := Responses[NoPermissionError.Key()]

		Design 1 was chosen because mapping between error type and error message is not necessarily 1-to-1.
		Multiple error types can produce the same error message. With design 1, it's also more obvious in the code
		to see which error is connected to which message key.
	*/
	switch err {
	case chatplugins.ErrCommandNotFound: // Command name is not found. Likely syncronization issue
		fallthrough
	case chatplugins.ErrNoPermission: // User has no permission
		fallthrough
	case chatplugins.ErrNotEnoughArguments: // Arguments
		fallthrough
	case chatplugins.ErrTargetCommandAlreadyExists: // Target command already exists and cannot be added
		fallthrough
	default:
		return models.DefaultFailureResponseKey
	}
}
