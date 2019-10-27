package chathandler

import (
	"log"
	"strings"

	"github.com/c-rainbow/simplechatbot/client"
	"github.com/c-rainbow/simplechatbot/models"
	"github.com/c-rainbow/simplechatbot/repository"

	pluginmanager "github.com/c-rainbow/simplechatbot/plugins/chat/manager"
	twitch_irc "github.com/gempir/go-twitch-irc"
)

type ChatMessageHandlerT interface {
	OnPrivateMessage(message twitch_irc.PrivateMessage)
}

type ChatMessageHandler struct {
	botInfo           *models.Bot
	repo              repository.SingleBotRepositoryT
	ircClient         client.TwitchClientT
	chatPluginManager pluginmanager.ChatCommandPluginManagerT
}

var _ ChatMessageHandlerT = (*ChatMessageHandler)(nil)

func NewChatMessageHandler(
	botInfo *models.Bot, repo repository.SingleBotRepositoryT, ircClient client.TwitchClientT,
	chatPluginManager pluginmanager.ChatCommandPluginManagerT) *ChatMessageHandler {
	return &ChatMessageHandler{botInfo: botInfo, repo: repo, ircClient: ircClient, chatPluginManager: chatPluginManager}
}

func (handler *ChatMessageHandler) OnPrivateMessage(message twitch_irc.PrivateMessage) {
	log.Println("Chat received: ", message.Raw)

	// TODO: Delete this hardcoded quit message.
	commandName := getCommandName(message.Message)
	commandName = strings.ToLower(commandName)
	if commandName == "!quit" && message.User.Name == "c_rainbow" {
		handler.ircClient.Depart(message.Channel)
		handler.ircClient.Disconnect()
	}

	handler.chatPluginManager.ProcessChat(message.Channel, &message.User, &message)
}

// Gets command name from the full chat text
func getCommandName(text string) string {
	// strings.Fields deals with heading/trailing/non-space whitespaces.
	fields := strings.Fields(text)
	return fields[0]
}
