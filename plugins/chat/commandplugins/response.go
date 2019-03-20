package commandplugins

import (
	"errors"

	"github.com/c-rainbow/simplechatbot/client"
	chatplugins "github.com/c-rainbow/simplechatbot/plugins/chat"
	"github.com/c-rainbow/simplechatbot/repository"

	models "github.com/c-rainbow/simplechatbot/models"
	parser "github.com/c-rainbow/simplechatbot/parser"
	twitch_irc "github.com/gempir/go-twitch-irc"
)

const (
	CommandResponsePluginType = "CommandResponsePluginType"
)

var (
	CommandNotFoundError   = errors.New("Command with the name is not found")
	NoDefaultResponseError = errors.New("Default response is not found for the command")
)

// Plugin that responds to user-defined chat commands.
type CommandResponsePlugin struct {
	ircClient client.TwitchClientT
	repo      repository.SingleBotRepositoryT
}

var _ chatplugins.ChatCommandPlugin = (*CommandResponsePlugin)(nil)

func NewCommandResponsePlugin(
	ircClient client.TwitchClientT, repo repository.SingleBotRepositoryT) chatplugins.ChatCommandPlugin {
	return &CommandResponsePlugin{ircClient: ircClient, repo: repo}
}

func (plugin *CommandResponsePlugin) Run(
	commandName string, channel string, sender *twitch_irc.User, message *twitch_irc.Message) error {
	// Read-action-print loop
	command, err := CommonRead(plugin.repo, commandName, channel, CommandResponsePluginType, sender, message)
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
func (plugin *CommandResponsePlugin) action(
	command *models.Command, channel string, sender *twitch_irc.User, message *twitch_irc.Message,
	err error) (string, error) {
	// Check what error is returned from Read()
	if err != nil {
		// TODO: Convert the error to public display message
		return err.Error(), nil
	}

	// Get default response
	response, exists := command.Responses[chatplugins.DefaultResponseKey]
	// This shouldn't happen in normal case. Default response always exists
	// if command is added in a correct way.
	if !exists {
		return "No response is found with the command", NoDefaultResponseError
	}

	// TODO: Support for API & function type variables
	converted, err := parser.ConvertResponse(&response, channel, sender, message)
	// This error usually happens when disabled/unsupported variable name is used.
	// It is usually already checked when command was added, but is double-checked here just in case.
	if err != nil {
		return "Failed to convert response to text", err
	}

	return converted, nil
}
