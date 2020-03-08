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
	Close()
}

// ChatMessageHandler chat message handler for a bot
type ChatMessageHandler struct {
	botInfo                 *models.Bot
	repo                    repository.SingleBotRepositoryT
	ircClient               client.TwitchClientT
	activePluginManager     pluginmanager.ActivePluginManagerT
	backgroundPluginManager pluginmanager.BackgroundPluginManagerT
}

var _ ChatMessageHandlerT = (*ChatMessageHandler)(nil)

// NewChatMessageHandler creates new handler
func NewChatMessageHandler(
	botInfo *models.Bot, repo repository.SingleBotRepositoryT, ircClient client.TwitchClientT,
	activePluginManager pluginmanager.ActivePluginManagerT,
	backgroundPluginManager pluginmanager.BackgroundPluginManagerT) ChatMessageHandlerT {
	return &ChatMessageHandler{botInfo: botInfo, repo: repo, ircClient: ircClient,
		activePluginManager: activePluginManager, backgroundPluginManager: backgroundPluginManager}
}

// DefaultChatMessageHandler creates new default handler
func DefaultChatMessageHandler(
	botInfo *models.Bot, repo repository.SingleBotRepositoryT, ircClient client.TwitchClientT) ChatMessageHandlerT {
	activePluginManager := pluginmanager.DefaultActivePluginManager(ircClient, repo)
	backgroundPluginManager := pluginmanager.NewBackgroundPluginManager(repo)
	return NewChatMessageHandler(botInfo, repo, ircClient, activePluginManager, backgroundPluginManager)
}

// OnPrivateMessage handles PRIVMSG
func (handler *ChatMessageHandler) OnPrivateMessage(message twitch_irc.PrivateMessage) {
	log.Println("Chat received: ", message.Raw)

	// TODO: Delete this hardcoded quit message.
	commandName := getCommandName(message.Message)
	commandName = strings.ToLower(commandName)
	if commandName == "!quit" && message.User.Name == "c_rainbow" {
		handler.ircClient.Depart(message.Channel)
		handler.ircClient.Disconnect()
	}

	handler.backgroundPluginManager.ProcessChat(message.Channel, &message.User, &message)
	handler.activePluginManager.ProcessChat(message.Channel, &message.User, &message)
}

// Close closes the chat message handler
func (handler *ChatMessageHandler) Close() {
	handler.backgroundPluginManager.Close()
	handler.activePluginManager.Close()
}

// Gets command name from the full chat text
func getCommandName(text string) string {
	// strings.Fields deals with heading/trailing/non-space whitespaces.
	fields := strings.Fields(text)
	return fields[0]
}
