package commandplugins

import (
	"sort"
	"strings"

	"github.com/c-rainbow/simplechatbot/client"
	chatplugins "github.com/c-rainbow/simplechatbot/plugins/chat"
	"github.com/c-rainbow/simplechatbot/repository"

	models "github.com/c-rainbow/simplechatbot/models"
	twitch_irc "github.com/gempir/go-twitch-irc"
)

const (
	ListCommandsPluginType = "ListCommandsPluginType"
)

type ListCommandsPluginFactory struct {
	ircClient client.TwitchClientT
	repo      repository.SingleBotRepositoryT
}

var _ chatplugins.ChatCommandPluginFactoryT = (*ListCommandsPluginFactory)(nil)

func NewListCommandsPluginFactory(
	ircClient client.TwitchClientT, repo repository.SingleBotRepositoryT) chatplugins.ChatCommandPluginFactoryT {
	return &ListCommandsPluginFactory{ircClient: ircClient, repo: repo}
}

func (plugin *ListCommandsPluginFactory) GetPluginType() string {
	return ListCommandsPluginType
}

func (plugin *ListCommandsPluginFactory) BuildNewPlugin() chatplugins.ChatCommandPluginT {
	return NewListCommandsPlugin(plugin.ircClient, plugin.repo)
}

// Plugin that responds to user-defined chat commands.
type ListCommandsPlugin struct {
	ircClient client.TwitchClientT
	repo      repository.SingleBotRepositoryT
}

var _ chatplugins.ChatCommandPluginT = (*ListCommandsPlugin)(nil)

func NewListCommandsPlugin(
	ircClient client.TwitchClientT, repo repository.SingleBotRepositoryT) chatplugins.ChatCommandPluginT {
	return &ListCommandsPlugin{ircClient: ircClient, repo: repo}
}

func (plugin *ListCommandsPlugin) GetPluginType() string {
	return ListCommandsPluginType
}

func (plugin *ListCommandsPlugin) ReactToChat(
	command *models.Command, channel string, sender *twitch_irc.User, message *twitch_irc.Message) {
	var targetCommands []*models.Command
	err := ValidateBasicInputs(command, channel, ListCommandsPluginType, sender, message)
	if err == nil {
		targetCommands, err = plugin.repo.ListCommands(channel)
	}

	responseText, err := plugin.GetResponseText(command, targetCommands, channel, sender, message, err)

	SendToChatClient(plugin.ircClient, channel, responseText)
	HandleError(err)
}

// Get response text of the executed command, based on the errors and progress so far.
func (plugin *ListCommandsPlugin) GetResponseText(
	command *models.Command, targetCommands []*models.Command, channel string, sender *twitch_irc.User,
	message *twitch_irc.Message, err error) (string, error) {
	// Get response key, build args, get parsed response, and convert it to text
	responseKey := plugin.GetResponseKey(err)
	args := []string{plugin.GetSortedCommandNames(targetCommands)}
	return ConvertToResponseText(command, responseKey, channel, sender, message, args)
}

func (plugin *ListCommandsPlugin) GetResponseKey(err error) string {
	// Normal case.
	if err == nil {
		return models.DefaultResponseKey
	}

	switch err {
	case chatplugins.ErrCommandNotFound: // Command name is not found. Likely synchronization issue
		fallthrough
	case chatplugins.ErrNoPermission: // User has no permission. Unlikely for this ListCommands plugin.
		fallthrough
	default:
		return models.DefaultFailureResponseKey
	}
}

func (plugin *ListCommandsPlugin) GetSortedCommandNames(commands []*models.Command) string {
	commandNames := make([]string, len(commands))
	for i, command := range commands {
		commandNames[i] = command.Name
	}

	sort.Strings(commandNames)
	return strings.Join(commandNames, ", ")
}
