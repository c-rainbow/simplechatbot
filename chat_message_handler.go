package simplechatbot

import (
	"strings"

	models "github.com/c-rainbow/simplechatbot/models"
	twitch_irc "github.com/gempir/go-twitch-irc"
)

type ChatMessageHandler struct {
	botInfo   *models.Bot
	repo      *SingleBotRepository
	ircClient *TwitchClient
}

func NewChatMessageHandler(
	botInfo *models.Bot, repo *SingleBotRepository, ircClient *TwitchClient) *ChatMessageHandler {
	return &ChatMessageHandler{botInfo: botInfo, repo: repo, ircClient: ircClient}
}

func (handler *ChatMessageHandler) OnNewMessage(channel string, sender twitch_irc.User, message twitch_irc.Message) {
	commandName := getCommandName(message.Text)
	// TODO: permission check, spam check, etc.
	command := handler.repo.GetCommandByChannelAndName(channel, commandName)
	if command != nil { //
		// response := command.Response
		// TODO: Format response and send
		// handler.ircClient.Say(response)
	}
}

// Gets command name from the full chat text
func getCommandName(text string) string {
	index := strings.Index(text, " ")
	// If there is no space in the chat text, then the chat itself is the command
	if index == -1 {
		return text
	}
	return text[:index]
}
