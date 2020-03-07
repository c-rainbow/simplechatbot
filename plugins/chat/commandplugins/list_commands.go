package commandplugins

import (
	"sort"
	"strings"

	"github.com/c-rainbow/simplechatbot/client"
	chatplugins "github.com/c-rainbow/simplechatbot/plugins/chat"
	"github.com/c-rainbow/simplechatbot/plugins/chat/common"
	"github.com/c-rainbow/simplechatbot/repository"

	models "github.com/c-rainbow/simplechatbot/models"
	twitch_irc "github.com/gempir/go-twitch-irc"
)

const (
	// ListCommandsPluginType plugin type name to list existing commands
	ListCommandsPluginType = "ListCommandsPluginType"
)

// ListCommandsPlugin plugin to list all chat commands
type ListCommandsPlugin struct {
	ircClient client.TwitchClientT
	repo      repository.SingleBotRepositoryT
}

var _ chatplugins.ChatCommandPluginT = (*ListCommandsPlugin)(nil)

// NewListCommandsPlugin creates a new ListCommandsPlugin
func NewListCommandsPlugin(
	ircClient client.TwitchClientT, repo repository.SingleBotRepositoryT) chatplugins.ChatCommandPluginT {
	return &ListCommandsPlugin{ircClient: ircClient, repo: repo}
}

// GetPluginType gets plugin type
func (plugin *ListCommandsPlugin) GetPluginType() string {
	return ListCommandsPluginType
}

// ReactToChat reacts to chat
func (plugin *ListCommandsPlugin) ReactToChat(
	command *models.Command, channel string, sender *twitch_irc.User, message *twitch_irc.PrivateMessage) {
	var targetCommands []*models.Command
	err := common.ValidateBasicInputs(command, channel, ListCommandsPluginType, sender, message)
	if err == nil {
		targetCommands, err = plugin.repo.ListCommands(channel)
	}

	responseText, err := plugin.GetResponseText(command, targetCommands, channel, sender, message, err)

	common.SendToChatClient(plugin.ircClient, channel, responseText)
	common.HandleError(err)
}

// GetResponseText gets response text of the executed command, based on the errors and progress so far.
func (plugin *ListCommandsPlugin) GetResponseText(
	command *models.Command, targetCommands []*models.Command, channel string, sender *twitch_irc.User,
	message *twitch_irc.PrivateMessage, err error) (string, error) {
	// Get response key, build args, get parsed response, and convert it to text
	responseKey := plugin.GetResponseKey(err)
	args := []string{plugin.getSortedCommandNames(targetCommands)}
	return common.ConvertToResponseText(command, responseKey, channel, sender, message, args)
}

// GetResponseKey returns response key from error type to build response text accordingly.
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

func (plugin *ListCommandsPlugin) getSortedCommandNames(commands []*models.Command) string {
	commandNames := make([]string, len(commands))
	for i, command := range commands {
		commandNames[i] = command.Name
	}

	sort.Strings(commandNames)
	return strings.Join(commandNames, ", ")
}
