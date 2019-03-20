package commandplugins

import (
	"errors"
	"log"

	client "github.com/c-rainbow/simplechatbot/client"
	models "github.com/c-rainbow/simplechatbot/models"
	parser "github.com/c-rainbow/simplechatbot/parser"
	plugins "github.com/c-rainbow/simplechatbot/plugins"
	chatplugins "github.com/c-rainbow/simplechatbot/plugins/chat"
	repository "github.com/c-rainbow/simplechatbot/repository"
	twitch_irc "github.com/gempir/go-twitch-irc"
)

var (
	NoPermissionError = errors.New("User has no permission to call this command")
)

// Common read function by add/edit/delete/respond commands.
// They are all identical except plugin type check part and actionFunction parameter
func CommonRun(repo repository.SingleBotRepositoryT, ircClient client.TwitchClientT, expectedPluginType string,
	actionFunction func(string, *models.Command) error, commandName string, channel string,
	sender *twitch_irc.User, message *twitch_irc.Message) error {

	// Read-action-print loop
	command, err := CommonRead(repo, commandName, channel, expectedPluginType, sender, message)
	toSay, err := CommonAction(repo.GetBotInfo(), actionFunction, command, channel, sender, message, err)
	return CommonOutput(ircClient, channel, toSay, err)
}

// Common read function by add/edit/delete/respond commands.
// They are all identical except plugin type check part
func CommonRead(repo repository.SingleBotRepositoryT, commandName string, channel string, expectedPluginType string,
	sender *twitch_irc.User, message *twitch_irc.Message) (*models.Command, error) {
	// Get command model from the repository
	command := repo.GetCommandByChannelAndName(channel, commandName)
	if command == nil {
		return nil, CommandNotFoundError
	}

	// Make sure that the plugin type is correct. This check was added for validation purposes,
	// because mismatch means data inconsistency between repository and chat message handler.
	if command.PluginType != expectedPluginType {
		return nil, plugins.IncorrectPluginTypeError
	}

	// Check permission. Still returns command object because permission error is user's misbehavior, and
	// not a plugin error. TODO: I am not 100% sure if "no permission" should be an error or return value.
	if err := UserHasPermission(channel, command, sender, message); err != nil {
		return command, err
	}

	return command, nil
}

// In this function, returned error means internal runtime error, not user input error.
// For example, NoPermissionsError is not an error here. However, a connection error to DB
// should be returned as error in this function.
//
// Note that CommandNotFoundError is also treated as an error, because in usual workflow,
// this plugin is only called after chat message handler checks for command name.
func CommonAction(
	botInfo *models.Bot, actionFunction func(string, *models.Command) error, command *models.Command,
	channel string, sender *twitch_irc.User, message *twitch_irc.Message, err error) (string, error) {
	// Check what error is returned from Read()
	if err != nil {
		// TODO: Convert the error to public display message
		return err.Error(), nil
	}

	// Parse command response
	parsedCommand, err := ParseCommand(botInfo.TwitchID, message.Text, channel, sender, message)
	if err != nil {
		return "Failed to parse command", err
	}

	// TODO:Validate

	// actionFunction is One of AddCommand, EditCommand, DeleteCommand
	err = actionFunction(channel, parsedCommand)
	if err != nil {
		// TODO: pull failure message for
		// failureResponse, exists := command.Responses[AddFailureKey]
		return "Failed to do actionFunction", err
	}

	// TODO: Add new variable, for the name of command just added
	// Get default response for successful add.
	successResponse, exists := command.Responses[chatplugins.DefaultResponseKey]
	// This shouldn't happen in normal case. Default response always exists
	// if command is added in a correct way.
	if !exists {
		return "No success response is found with the command", nil
	}

	// TODO: Support for API & function type variables
	converted, err := parser.ConvertResponse(&successResponse, channel, sender, message)
	// This error usually happens when disabled/unsupported variable name is used.
	// It is usually already checked when command was added, but is double-checked here just in case.
	if err != nil {
		return "Failed to convert response to text", err
	}

	return converted, nil
}

func CommonOutput(ircClient client.TwitchClientT, channel string, toSay string, err error) error {
	// Even with error, the plugin might respond that "There is unknown error"
	if toSay != "" {
		ircClient.Say(channel, toSay)
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

// TODO: Why does this need boolean value?
func UserHasPermission(channel string, command *models.Command, sender *twitch_irc.User,
	message *twitch_irc.Message) error {
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
		return nil
	}*/
	return nil
}
