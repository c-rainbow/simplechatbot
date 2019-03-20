package commandplugins

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	models "github.com/c-rainbow/simplechatbot/models"
	parser "github.com/c-rainbow/simplechatbot/parser"
	chatplugins "github.com/c-rainbow/simplechatbot/plugins/chat"
	twitch_irc "github.com/gempir/go-twitch-irc"
)

var (
	NotEnoughArgumentsError = errors.New("Not enough arguments")
)

func ParseCommand(
	botID int64, text string, channel string, sender *twitch_irc.User,
	message *twitch_irc.Message) (*models.Command, error) {

	name, response, err := ParseNameAndResponseFromChat(text)
	if err != nil {
		return nil, err
	}

	channelID, err := strconv.Atoi(message.ChannelID)

	parsedResponse := parser.ParseResponse(response)
	responseMap := make(map[string]parser.ParsedResponse)
	responseMap[chatplugins.DefaultResponseKey] = *parsedResponse
	fmt.Println("responseMap: ", responseMap)

	// Parse command and response
	command := models.Command{
		BotID:          botID,
		ChannelID:      int64(channelID),
		Name:           name,
		PluginType:     CommandResponsePluginType,
		Responses:      responseMap,
		CooldownSecond: 5,
		Enabled:        true,
		Group:          "",
	}
	return &command, nil
}

func ParseNameAndResponseFromChat(text string) (string, string, error) {
	// TODO: This does not handle consecutive whitespaces in response text.
	fields := strings.Fields(text)

	// This method is called only when !addcom/!editcom/!delcom is called from chat.
	// Minimum length 3 is correct (!*com [commandName] [response]) for addcom/editcom
	// For delcom, length should be 2
	if len(fields) < 2 {
		return "", "", NotEnoughArgumentsError
	}

	response := ""
	if len(fields) < 3 {
		return fields[1], "", nil
	} else {
		response = strings.Join(fields[2:], " ")
	}

	return fields[1], response, nil
}
