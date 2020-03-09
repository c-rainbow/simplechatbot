package commandplugins

import (
	"github.com/c-rainbow/simplechatbot/client"
	chatplugins "github.com/c-rainbow/simplechatbot/plugins/chat"
	"github.com/c-rainbow/simplechatbot/plugins/chat/common"
	"github.com/c-rainbow/simplechatbot/repository"

	"github.com/c-rainbow/simplechatbot/models"
	twitch_irc "github.com/gempir/go-twitch-irc"
)

const (
	// CommandResponsePluginType plugin type name to respond to a chat command
	CommandResponsePluginType = "CommandResponsePluginType"
)

// CommandResponsePlugin plugin that responds to user-defined chat commands.
type CommandResponsePlugin struct {
	ircClient client.TwitchClientT
	repo      repository.SingleBotRepositoryT
}

var _ chatplugins.ChatCommandPluginT = (*CommandResponsePlugin)(nil)

// NewCommandResponsePlugin creates a new CommandResponsePlugin
func NewCommandResponsePlugin(
	ircClient client.TwitchClientT, repo repository.SingleBotRepositoryT) chatplugins.ChatCommandPluginT {
	return &CommandResponsePlugin{ircClient: ircClient, repo: repo}
}

// GetPluginType gets plugin type
func (plugin *CommandResponsePlugin) GetPluginType() string {
	return CommandResponsePluginType
}

// ReactToChat reacts to chat
func (plugin *CommandResponsePlugin) ReactToChat(
	command *models.Command, channel string, sender *twitch_irc.User, message *twitch_irc.PrivateMessage) {
	// Basic validations
	err := common.ValidateBasicInputs(command, channel, CommandResponsePluginType, sender, message)

	// This should get response object, resolve variables, then get response text from it
	responseText, err := plugin.GetResponseText(command, channel, sender, message, err)

	common.SendToChatClient(plugin.ircClient, channel, responseText)
	common.HandleError(err)
}

// GetResponseText gets response text of the executed command, based on the errors and progress so far.
func (plugin *CommandResponsePlugin) GetResponseText(
	command *models.Command, channel string, sender *twitch_irc.User, message *twitch_irc.PrivateMessage,
	err error) (string, error) {
	// Get response key, build args, get parsed response, and convert it to text
	responseKey := plugin.GetResponseKey(err)
	args := []string{}
	return common.ConvertToResponseText(command, responseKey, channel, sender, message, args)
}

// GetResponseKey returns response key from error type to build response text accordingly.
func (plugin *CommandResponsePlugin) GetResponseKey(err error) string {
	// Normal case.
	if err == nil {
		return models.DefaultResponseKey
	}

	switch err {
	case chatplugins.ErrCommandNotFound: // Command name is not found. Likely synchronization issue
		fallthrough
	case chatplugins.ErrNoPermission: // User has no permission. Unlikely for this CommandResponse plugin.
		fallthrough
	default:
		return models.DefaultFailureResponseKey
	}
}
