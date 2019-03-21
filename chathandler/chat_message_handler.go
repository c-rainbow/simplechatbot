package chathandler

import (
	"fmt"
	"strings"

	"github.com/c-rainbow/simplechatbot/client"
	"github.com/c-rainbow/simplechatbot/models"
	"github.com/c-rainbow/simplechatbot/repository"

	pluginmanager "github.com/c-rainbow/simplechatbot/plugins/chat/manager"
	twitch_irc "github.com/gempir/go-twitch-irc"
)

type ChatMessageHandlerT interface {
	OnNewMessage(channel string, sender twitch_irc.User, message twitch_irc.Message)
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

func (handler *ChatMessageHandler) OnNewMessage(
	channel string, sender twitch_irc.User, message twitch_irc.Message) {
	fmt.Println("Chat received: ", message.Raw)

	// TODO: Delete this hardcoded quit message.
	commandName := getCommandName(message.Text)
	commandName = strings.ToLower(commandName)
	if commandName == "!quit" && sender.Username == "c_rainbow" {
		handler.ircClient.Depart(channel)
		handler.ircClient.Disconnect()
	}

	handler.chatPluginManager.ProcessChat(channel, &sender, &message)
	// TODO: If command is retrieved, why not use this directly?
	// Check if command with the same name exists
	// The repository (of type SingleBotRepositoryT) ensures that the correct bot responds to the command.
	// If this line is executed in a different bot, nil would be returned.
	/*command := handler.repo.GetCommandByChannelAndName(channel, commandName)
	if command == nil { // Chat is not a bot command.
		return
	}

	msgQueue := make(chan *models.Command)

	msgQueue <- command

	var plugin pluginmanager.ChatCommandPluginT
	// TODO: Eventually, get plugin from a factory or pluginmanager
	switch command.PluginType {
	// Simply responding to command, no other action
	case commandplugins.CommandResponsePluginType:
		plugin = commandplugins.NewCommandResponsePlugin(handler.ircClient, handler.repo)
	// Add command, like !addcom
	case commandplugins.AddCommandPluginType:
		plugin = commandplugins.NewAddCommandPlugin(handler.ircClient, handler.repo)
	// Edit command, like !editcom
	case commandplugins.EditCommandPluginType:
		plugin = commandplugins.NewEditCommandPlugin(handler.ircClient, handler.repo)
	// Delete command, like !delcom
	case commandplugins.DeleteCommandPluginType:
		plugin = commandplugins.NewDeleteCommandPlugin(handler.ircClient, handler.repo)
	// List commands, like !commands
	case commandplugins.ListCommandsPluginType:
		plugin = commandplugins.NewListCommandsPlugin(handler.ircClient, handler.repo)
	default:
		log.Println("Unknown plugin type '", command.PluginType, "'")
	}

	if plugin != nil {
		go func() {

			err := plugin.Run(commandName, channel, &sender, &message)
			if err != nil {
				log.Println("Error while running plugin for '", commandName, "': ", err.Error())
			}
		}()
	}*/

}

// Gets command name from the full chat text
func getCommandName(text string) string {
	// strings.Fields deals with heading/trailing/non-space whitespaces.
	fields := strings.Fields(text)
	return fields[0]
}
