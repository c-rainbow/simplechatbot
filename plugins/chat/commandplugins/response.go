package commandplugins

import (
	"github.com/c-rainbow/simplechatbot/client"
	"github.com/c-rainbow/simplechatbot/plugins/chat"
	"github.com/c-rainbow/simplechatbot/repository"

	"github.com/c-rainbow/simplechatbot/models"
	twitch_irc "github.com/gempir/go-twitch-irc"
)

const (
	CommandResponsePluginType = "CommandResponsePluginType"
)

type CommandResponsePluginFactory struct {
	ircClient client.TwitchClientT
	repo      repository.SingleBotRepositoryT
}

var _ chatplugins.ChatCommandPluginFactoryT = (*CommandResponsePluginFactory)(nil)

func NewCommandResponsePluginFactory(
	ircClient client.TwitchClientT, repo repository.SingleBotRepositoryT) chatplugins.ChatCommandPluginFactoryT {
	return &CommandResponsePluginFactory{ircClient: ircClient, repo: repo}
}

func (plugin *CommandResponsePluginFactory) GetPluginType() string {
	return CommandResponsePluginType
}

func (plugin *CommandResponsePluginFactory) BuildNewPlugin() chatplugins.ChatCommandPluginT {
	return NewCommandResponsePlugin(plugin.ircClient, plugin.repo)
}

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
	command *models.Command, channel string, sender *twitch_irc.User, message *twitch_irc.Message) {
		// Basic validations
	err := ValidateBasicInputs(command, channel, CommandResponsePluginType, sender, message)

	responseText, err := plugin.GetResponseText(command, channel, sender, message, err)

	SendToChatClient(plugin.ircClient, channel, responseText)
	HandleError(err)
}

func (plugin *CommandResponsePlugin) GetResponseText(
	command *models.Command, channel string, sender *twitch_irc.User, message *twitch_irc.Message,
	err error) (string, error) {
	// Get response key, build args, get parsed response, and convert it to text
	responseKey := plugin.GetResponseKey(err)
	args := []string{}
	return ConvertToResponseText(command, responseKey, channel, sender, message, args)
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
