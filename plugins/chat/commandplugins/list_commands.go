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

func (plugin *ListCommandsPlugin) Run(
	command *models.Command, channel string, sender *twitch_irc.User, message *twitch_irc.Message) error {
	// Read-action-print loop
	command, err := CommonRead(plugin.repo, command, channel, ListCommandsPluginType, sender, message)
	toSay, err := plugin.action(command, channel, sender, message, err)
	err = CommonOutput(plugin.ircClient, channel, toSay, err)
	if err != nil {
		return err
	}
	return nil
}

// In this function, returned error means internal runtime error, not user input error.
// For example, NoPermissionsError is not an error here. However, a connection error to DB
// should be returned as error in this function.
//
// Note that CommandNotFoundError is also treated as an error, because in usual workflow,
// this plugin is only called after chat message handler checks for command name.
func (plugin *ListCommandsPlugin) action(
	command *models.Command, channel string, sender *twitch_irc.User, message *twitch_irc.Message,
	err error) (string, error) {
	// Check what error is returned from Read()
	if err != nil {
		// TODO: Convert the error to public display message
		return err.Error(), nil
	}

	// Get all command names, sort them, and join by comma
	commands, err := plugin.repo.ListCommands(channel)
	commandNames := make([]string, len(commands))
	for i, command := range commands {
		commandNames[i] = command.Name
	}

	sort.Strings(commandNames)
	return strings.Join(commandNames, ", "), nil
}

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
