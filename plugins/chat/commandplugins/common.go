package commandplugins

import (
	"log"

	client "github.com/c-rainbow/simplechatbot/client"
	models "github.com/c-rainbow/simplechatbot/models"
	parser "github.com/c-rainbow/simplechatbot/parser"
	plugins "github.com/c-rainbow/simplechatbot/plugins"
	chatplugins "github.com/c-rainbow/simplechatbot/plugins/chat"
	repository "github.com/c-rainbow/simplechatbot/repository"
	twitch_irc "github.com/gempir/go-twitch-irc"
)

// Common read function by add/edit/delete/respond commands.
// They are all identical except plugin type check part and actionFunction parameter
func CommonRun(repo repository.SingleBotRepositoryT, ircClient client.TwitchClientT, expectedPluginType string,
	actionFunction func(string, *models.Command) error, command *models.Command, channel string,
	sender *twitch_irc.User, message *twitch_irc.Message) error {

	// Read-action-print loop
	err := CommonValidateInputs(command, channel, expectedPluginType, sender, message)
	//return CommonOutput(ircClient, channel, toSay, err)
	return nil
}

// Validate basic inputs. Check for existence, plugin type, and user permission of the command
func CommonValidateInputs(command *models.Command, channel string, expectedPluginType string, sender *twitch_irc.User,
	message *twitch_irc.Message) error {
	if command == nil {
		return chatplugins.ErrCommandNotFound
	}

	// Make sure that the plugin type is correct. This check was added for validation purposes,
	// because mismatch means data inconsistency between repository and chat message handler.
	if command.PluginType != expectedPluginType {
		return plugins.ErrIncorrectPluginType
	}

	// Check permission
	if !UserHasPermission(channel, command, sender, message) {
		return chatplugins.ErrNoPermission
	}

	return nil
}

/*
Common action performed by AddCommand/EditCommand.

In this function, returned error means internal runtime error, not user input error.
For example, NoPermissionsError is not an error here. However, a connection error to DB
should be returned as error in this function.

Note that CommandNotFoundError is also treated as an error, because in usual workflow,
this plugin is only called after chat message handler checks for command name.

Returned string is the message to send to the output method (which is then sent to the IRC client)
If error is not
*/
/*func CommonAction(
	botInfo *models.Bot, actionFn func(string, *models.Command) error, command *models.Command, channel string,
	sender *twitch_irc.User, message *twitch_irc.Message) (string, error) {

	// Parse chat message to command struct
	parsedCommand, err := ParseCommand(botInfo.TwitchID, message.Text, channel, sender, message)
	if err != nil {
		return "Failed to parse chat message to command", err
	}

	// Validate responses..
	for _, response := range parsedCommand.Responses {
		err = parser.Validate(&response)
		if err != nil {
			return "Response text cannot be validated", err
		}
	}

	// actionFunction is One of AddCommand or EditCommand
	err = actionFn(channel, parsedCommand)
	if err != nil {
		return "Error returned from actionFunction", err
	}

	args := []string{parsedCommand.Name}
	converted, err := CommonConvertToResponseText(command, channel, sender, message, args)

	return converted, nil
}*/

func CommonConvertToResponseText(
	command *models.Command, responseKey string, channel string, sender *twitch_irc.User, message *twitch_irc.Message,
	args []string) (string, error) {

	// Get parsed response with the response key
	parsedResponse, exists := command.Responses[responseKey]
	if !exists {
		// Default response should always exists. This is software bug.
		log.Println("Response text cannot be found for", command.Name, "in", channel)
		return "", chatplugins.ErrResponseNotFound
	}

	// Try parsing the response with parameters
	converted, err := parser.ConvertResponse(&parsedResponse, channel, sender, message, args)
	// This error usually happens when disabled/unsupported variable name is used.
	// It is usually already checked when command was added, but is double-checked here just in case.
	if err != nil {
		log.Println("Failed to convert response to text", err.Error())
		return "", err
	}

	return converted, nil
}

// Check if user has permission for a command. Checking follower status requires call to Twitch API,
// so it is not included here.
// TODO: How to check VIP and staff status?
func UserHasPermission(channel string, command *models.Command, sender *twitch_irc.User,
	message *twitch_irc.Message) bool {
	// (1) Broadcaster can do everything
	if sender.Username == channel {
		return true
	}

	// (2) If permission is everyone, return true
	if command.Permission == models.PermissionEveryone {
		return true
	}

	// (3) If permission is subscriber, and user is subscriber
	isSubscriber := message.Tags["subscriber"]
	if isSubscriber == "1" && (command.Permission&models.PermissionSubscriber) > 0 {
		return true
	}

	// (4) If permission is moderator, and user is moderator
	isMod := message.Tags["mod"]
	if isMod == "1" && (command.Permission&models.PermissionModerator) > 0 {
		return true
	}

	return false
}

// Common function used by all chat command plugins to output to the IRC channel
func CommonOutput(ircClient client.TwitchClientT, channel string, toSay string) {
	if toSay != "" {
		ircClient.Say(channel, toSay)
	}
}

// TODO: This function may be used for sending error for analysis, etc.
func CommonHandleError(err error) {
	// Empty for now
	return
}
