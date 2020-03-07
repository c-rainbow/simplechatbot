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

// ChatMessageHandlerT is the general interface for chat message handler
type ChatMessageHandlerT interface {
	OnPrivateMessage(message twitch_irc.PrivateMessage)
}

// ChatMessageHandler chat message handler for a bot
type ChatMessageHandler struct {
	botInfo                  *models.Bot
	repo                     repository.SingleBotRepositoryT
	ircClient                client.TwitchClientT
	chatCommandPluginManager pluginmanager.ChatCommandPluginManagerT
}

var _ ChatMessageHandlerT = (*ChatMessageHandler)(nil)

// NewChatMessageHandler creates new handler
func NewChatMessageHandler(
	botInfo *models.Bot, repo repository.SingleBotRepositoryT, ircClient client.TwitchClientT,
	chatCommandPluginManager pluginmanager.ChatCommandPluginManagerT) *ChatMessageHandler {
	return &ChatMessageHandler{botInfo: botInfo, repo: repo, ircClient: ircClient,
		chatCommandPluginManager: chatCommandPluginManager}
}

// OnPrivateMessage handles PRIVMSG
func (handler *ChatMessageHandler) OnPrivateMessage(message twitch_irc.PrivateMessage) {
	log.Println("Chat received: ", message.Raw)

	// TODO: Do the entire process in a separate goroutine
	// This OnPrivate Message may take a bit longer when spamcheck & etc are added.

	// TODO: Delete this hardcoded quit message.
	commandName := getCommandName(message.Message)
	commandName = strings.ToLower(commandName)
	if commandName == "!quit" && message.User.Name == "c_rainbow" {
		handler.ircClient.Depart(message.Channel)
		handler.ircClient.Disconnect()
	}

	handler.chatCommandPluginManager.ProcessChat(message.Channel, &message.User, &message)
}

// Gets command name from the full chat text
func getCommandName(text string) string {
	// strings.Fields deals with heading/trailing/non-space whitespaces.
	fields := strings.Fields(text)
	return fields[0]
}
