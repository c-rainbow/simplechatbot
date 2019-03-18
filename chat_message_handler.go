package simplechatbot

import (
	"fmt"
	"log"
	"strings"

	"github.com/c-rainbow/simplechatbot/commands"
	models "github.com/c-rainbow/simplechatbot/models"
	twitch_irc "github.com/gempir/go-twitch-irc"
)

type ChatMessageHandlerT interface {
	OnNewMessage(channel string, sender twitch_irc.User, message twitch_irc.Message)
}

type ChatMessageHandler struct {
	botInfo   *models.Bot
	repo      SingleBotRepositoryT
	ircClient *TwitchClient
	chatters  map[string]bool
}

var _ ChatMessageHandlerT = (*ChatMessageHandler)(nil)

func NewChatMessageHandler(
	botInfo *models.Bot, repo SingleBotRepositoryT, ircClient *TwitchClient) *ChatMessageHandler {
	return &ChatMessageHandler{botInfo: botInfo, repo: repo, ircClient: ircClient, chatters: make(map[string]bool)}
}

func (handler *ChatMessageHandler) OnNewMessage(
	channel string, sender twitch_irc.User, message twitch_irc.Message) {
	fmt.Println("Chat received: ", message.Raw)
	commandName := getCommandName(message.Text)
	commandName = strings.ToLower(commandName)

	// One of reserved command names

	if commandName == "!addcom" {
		// Get
		parsedCommand, err := commands.ParseCommand(handler.botInfo.TwitchID, message.Text, channel, &sender, &message)
		if err != nil {
			log.Println("!addcom err: ", err.Error())
			return
		}
		handler.repo.AddCommand(channel, parsedCommand)

	} else if commandName == "!editcom" {

	} else if commandName == "!delcom" {

	}

	// Check if command with the same name exists
	command := handler.repo.GetCommandByChannelAndName(channel, commandName)
	displayName := sender.DisplayName
	toSay := ""

	if commandName == "hello" || commandName == "hi" {
		toSay = "Hi " + displayName
	} else if commandName == "안녕하세요" && sender.Username != "c_rainbow" {
		toSay = displayName + " 님도 안녕하세요"
	} else if commandName == "!quit" && sender.Username == "c_rainbow" {
		handler.ircClient.Depart(channel)
		handler.ircClient.Disconnect()
	}

	// Check for new chatter
	if _, has := handler.chatters[displayName]; !has {
		handler.chatters[displayName] = true
		toSay = displayName + " 님 어서오세요 환영합니다"
	}

	if toSay != "" {
		handler.ircClient.Say(channel, toSay)
	}

	// TODO: permission check, spam check, etc.
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
