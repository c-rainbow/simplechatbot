package simplechatbot

import (
	"fmt"
	"strings"

	twitch_irc "github.com/gempir/go-twitch-irc"
)

type ChatMessageHandler struct {
	// commands to handle
	userCommandsMap map[string]map[string]*Command
	botDataManager  *BotDataManager
	ircClient       *TwitchClient
}

func NewChatMessageHandler(ircClient *TwitchClient, botDataManager *BotDataManager) *ChatMessageHandler {
	userCommands := (*botDataManager).GetAllUserCommands()

	userCommandsMap := make(map[string]map[string]*Command)
	for _, userCommand := range userCommands {
		commandMap := make(map[string]*Command)
		userCommandsMap[userCommand.User.Username] = commandMap
		for _, command := range userCommand.Commands {
			commandMap[command.Name] = command
		}
	}

	return &ChatMessageHandler{
		userCommandsMap: userCommandsMap, botDataManager: botDataManager, ircClient: ircClient}
}

func (handler *ChatMessageHandler) OnNewMessage(channel string, user twitch_irc.User, message twitch_irc.Message) {
	fmt.Println("Chat received: " + message.Text)
	fmt.Println("Chat received raw: " + message.Raw)
	commandText := strings.SplitN(message.Text, " ", 2)
	if commandMap, ok := handler.userCommandsMap[channel]; ok {
		if command, ok := commandMap[commandText[0]]; ok {
			handler.ircClient.Say(channel, command.Response)
		}
	}

}
