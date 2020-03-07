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
	CommandResponsePluginType = "CommandResponsePluginType"
)

// Plugin that responds to user-defined chat commands.
type CommandResponsePlugin struct {
	ircClient client.TwitchClientT
	repo      repository.SingleBotRepositoryT
}

var _ chatplugins.ChatCommandPluginT = (*CommandResponsePlugin)(nil)

func NewCommandResponsePlugin(
	ircClient client.TwitchClientT, repo repository.SingleBotRepositoryT) chatplugins.ChatCommandPluginT {
	return &CommandResponsePlugin{ircClient: ircClient, repo: repo}
}

func (plugin *CommandResponsePlugin) GetPluginType() string {
	return CommandResponsePluginType
}

func (plugin *CommandResponsePlugin) ReactToChat(
	command *models.Command, channel string, sender *twitch_irc.User, message *twitch_irc.PrivateMessage) {
	// Basic validations
	err := common.ValidateBasicInputs(command, channel, CommandResponsePluginType, sender, message)

	responseText, err := plugin.GetResponseText(command, channel, sender, message, err)

	common.SendToChatClient(plugin.ircClient, channel, responseText)
	common.HandleError(err)
}

func (plugin *CommandResponsePlugin) GetResponseText(
	command *models.Command, channel string, sender *twitch_irc.User, message *twitch_irc.PrivateMessage,
	err error) (string, error) {
	// Get response key, build args, get parsed response, and convert it to text
	responseKey := plugin.GetResponseKey(err)
	args := []string{}
	return common.ConvertToResponseText(command, responseKey, channel, sender, message, args)
}

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
