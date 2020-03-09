package common

import (
	"log"

	"github.com/c-rainbow/simplechatbot/client"
	"github.com/c-rainbow/simplechatbot/models"
	"github.com/c-rainbow/simplechatbot/parser"
	"github.com/c-rainbow/simplechatbot/plugins"
	chatplugins "github.com/c-rainbow/simplechatbot/plugins/chat"
	twitch_irc "github.com/gempir/go-twitch-irc"
)

// Validate basic inputs. Check for existence, plugin type, and user permission of the command
func ValidateBasicInputs(
	command *models.Command, channel string, expectedPluginType string, sender *twitch_irc.User,
	message *twitch_irc.PrivateMessage) error {
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

func ConvertToResponseText(
	command *models.Command, responseKey string, channel string, sender *twitch_irc.User,
	message *twitch_irc.PrivateMessage, args []string) (string, error) {

	// Get parsed response with the response key
	parsedResponse, exists := command.Responses[responseKey]
	if !exists {
		// Default response should always exists. This is software bug.
		log.Println("Response text cannot be found for", command.Name, "in", channel)
		return "", chatplugins.ErrResponseNotFound
	}

	// Try parsing the response with parameters
	converted, err := parser.ConvertResponse(&parsedResponse, channel, sender, message, nil, args)
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
	message *twitch_irc.PrivateMessage) bool {
	// (1) Broadcaster can do everything
	if sender.Name == channel {
		return true
	}

	// (2) If permission is everyone, return true
	if command.Permission&models.PermissionEveryone > 0 {
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
func SendToChatClient(ircClient client.TwitchClientT, channel string, toSay string) {
	if toSay != "" {
		ircClient.Say(channel, toSay)
	}
}

// TODO: This function may be used for sending error for analysis, etc.
func HandleError(err error) {
	if err != nil {
		log.Println("Some Error: ", err.Error())
	}
	// Empty for now
	return
}
