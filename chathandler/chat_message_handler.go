package chathandler

import (
	"fmt"
	"log"
	"strings"

	client "github.com/c-rainbow/simplechatbot/client"
	models "github.com/c-rainbow/simplechatbot/models"
	repository "github.com/c-rainbow/simplechatbot/repository"

	chat_plugins "github.com/c-rainbow/simplechatbot/plugins/chat"
	commandplugins "github.com/c-rainbow/simplechatbot/plugins/chat/commandplugins"
	twitch_irc "github.com/gempir/go-twitch-irc"
)

type ChatMessageHandlerT interface {
	OnNewMessage(channel string, sender twitch_irc.User, message twitch_irc.Message)
}

type ChatMessageHandler struct {
	botInfo   *models.Bot
	repo      repository.SingleBotRepositoryT
	ircClient client.TwitchClientT
	chatters  map[string]bool
}

var _ ChatMessageHandlerT = (*ChatMessageHandler)(nil)

func NewChatMessageHandler(
	botInfo *models.Bot, repo repository.SingleBotRepositoryT, ircClient client.TwitchClientT) *ChatMessageHandler {
	return &ChatMessageHandler{botInfo: botInfo, repo: repo, ircClient: ircClient, chatters: make(map[string]bool)}
}

func (handler *ChatMessageHandler) OnNewMessage(
	channel string, sender twitch_irc.User, message twitch_irc.Message) {
	fmt.Println("Chat received: ", message.Raw)
	commandName := getCommandName(message.Text)
	commandName = strings.ToLower(commandName)

	// TODO: Delete this hardcoded quit message.
	if commandName == "!quit" && sender.Username == "c_rainbow" {
		handler.ircClient.Depart(channel)
		handler.ircClient.Disconnect()
	}

	// Check if command with the same name exists
	// The repository (of type SingleBotRepositoryT) ensures that the correct bot responds to the command.
	// If this line is executed in a different bot, nil would be returned.
	command := handler.repo.GetCommandByChannelAndName(channel, commandName)
	if command == nil { // Chat is not a bot command.
		return
	}

	var plugin chat_plugins.ChatCommandPlugin
	// TODO: Eventually, get plugin from a factory or manager
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
		err := plugin.Run(commandName, channel, &sender, &message)
		if err != nil {
			log.Println("Error while running plugin for '", commandName, "': ", err.Error())
		}
	}

}

// Gets command name from the full chat text
func getCommandName(text string) string {
	// strings.Fields deals with heading/trailing/non-space whitespaces.
	fields := strings.Fields(text)
	return fields[0]
}
