package commandplugins

import (
	"errors"
	"log"

	"github.com/c-rainbow/simplechatbot/commands"

	bot "github.com/c-rainbow/simplechatbot"
	models "github.com/c-rainbow/simplechatbot/models"
	parser "github.com/c-rainbow/simplechatbot/parser"
	twitch_irc "github.com/gempir/go-twitch-irc"
)

var (
	CommandNotFoundError   = errors.New("Command with the name is not found")
	NoPermissionError      = errors.New("User has no permission to call this command")
	NoDefaultResponseError = errors.New("Default response is not found for the command")
)

// Plugin that responds to user-defined chat commands.
type CommandResponsePlugin struct {
	ircClient *bot.TwitchClient
	repo      bot.SingleBotRepositoryT
}

var _ ChatCommandPlugin = (*CommandResponsePlugin)(nil)

func (plugin *CommandResponsePlugin) Run(
	commandName string, channel string, sender *twitch_irc.User, message *twitch_irc.Message) error {
	// Read-action-print loop
	command, err := plugin.read(commandName, channel, sender, message)
	toSay, err := plugin.action(command, channel, sender, message, err)
	err = plugin.output(channel, toSay, err)
	if err != nil {
		return err
	}
	return nil
}

func (plugin *CommandResponsePlugin) read(commandName string, channel string, sender *twitch_irc.User,
	message *twitch_irc.Message) (*models.Command, error) {
	// Get command model from the repository
	command := plugin.repo.GetCommandByChannelAndName(channel, commandName)
	if command == nil {
		return nil, CommandNotFoundError
	}

	// Check permission
	allowed, err := plugin.userHasPermission(channel, commandName, sender, message)
	if err != nil {
		return nil, err
	}
	if !allowed {
		return nil, NoPermissionError
	}

	return command, nil
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
	response, exists := command.Responses[commands.DefaultResponseKey]
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

func (plugin *CommandResponsePlugin) output(channel string, toSay string, err error) error {
	// Even with error, the plugin might respond that "There is unknown error"
	if toSay != "" {
		plugin.ircClient.Say(channel, toSay)
	}
	if err != nil {
		// Don't print anything because this is abnormal case.
		// In normal workflow, chat message handler already checks for existence
		// TODO: log error in more detail
		log.Fatal("Unexpected error:", err)
		return err
	}

	return nil
}

func (plugin *CommandResponsePlugin) userHasPermission(
	channel string, commandName string, sender *twitch_irc.User, message *twitch_irc.Message) (bool, error) {
	/* TODO:
	(1) Update function signature to accept user level (or is included in tags?)
	(2) If permission is everyone, return (true, nil)
	(3) If permission is follower, check follow status
	(4) If permission is subscriber, check subscriber status
	(5) If permission is vip, check vip status
	(6) If permission is moderator, check mod status
	(7) If permission is streamer, check streamer status

	TODO: How to check staff status?
	*/
	/*isMod := message.Tags["mod"]
	isBroadcaster := sender.Username == channel
	// TODO: How to check if user is staff?
	if isMod == "1" || isBroadcaster {
		return true, nil
	}*/
	return true, nil
}