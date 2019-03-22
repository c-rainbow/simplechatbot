package commandplugins

import (
	"strconv"
	"strings"

	models "github.com/c-rainbow/simplechatbot/models"
	parser "github.com/c-rainbow/simplechatbot/parser"
	chatplugins "github.com/c-rainbow/simplechatbot/plugins/chat"
	twitch_irc "github.com/gempir/go-twitch-irc"
)

var ()

func ParseCommand(
	botID int64, text string, channel string, sender *twitch_irc.User,
	message *twitch_irc.Message) (*models.Command, error) {

	name, response := GetCommandNameAndResponseTextFromChat(text)

	channelID, err := strconv.Atoi(message.ChannelID)
	if err != nil {
		return nil, err
	}

	parsedResponse := parser.ParseResponse(response)
	responseMap := make(map[string]models.ParsedResponse)
	responseMap[chatplugins.DefaultResponseKey] = *parsedResponse

	// Parse command and response
	command := models.Command{
		BotID:          botID,
		ChannelID:      int64(channelID),
		Name:           name,
		PluginType:     CommandResponsePluginType,
		Permission:     chatplugins.PermissionEveryone,
		Responses:      responseMap,
		CooldownSecond: 5,
		Enabled:        true,
		Group:          "",
	}
	return &command, nil
}

func GetTargetCommandNameAndResponse(text string) (string, string) {
	// TODO: This function does not acknowledge consecutive whitespaces in response text.
	// For example, if user types "!addcom !hello Welcome  \t  $(user)     here!", then
	// the response will be shortened to "Welcome $(user) here!", removing all long whitespaces
	// between words.
	fields := strings.Fields(text)

	// For DeleteCommand plugin
	if len(fields) == 2 {
		return fields[1], ""
	}

	// For AddCommand/EditCommand plugin
	response := strings.Join(fields[2:], " ")
	return fields[1], response

}
