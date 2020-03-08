package manager

import (
	"log"
	"strings"

	"github.com/c-rainbow/simplechatbot/plugins/chat/games"
	"github.com/c-rainbow/simplechatbot/plugins/chat/selfban"

	"github.com/c-rainbow/simplechatbot/client"
	chatplugins "github.com/c-rainbow/simplechatbot/plugins/chat"
	"github.com/c-rainbow/simplechatbot/plugins/chat/commandplugins"
	"github.com/c-rainbow/simplechatbot/repository"
	twitch_irc "github.com/gempir/go-twitch-irc"
)

// ChatCommandPluginManagerT interface for chat command plugin manager
type ChatCommandPluginManagerT interface {
	RegisterPlugin(plugin chatplugins.ChatCommandPluginT)
	ProcessChat(channel string, sender *twitch_irc.User, message *twitch_irc.PrivateMessage)
}

/*
ChatCommandPluginManager Chat command plugin manager. All chat command plugins are registered here,
and depends on the command, the manager distributes command struct to the right channel.

Each type of plugin has a channel and a worker pool with at least 1 worker, which keeps
polling from the channel until it closes.
*/
type ChatCommandPluginManager struct {
	repo      repository.SingleBotRepositoryT
	pluginMap map[string]chatplugins.ChatCommandPluginT // From plugin type name to plugin object
}

var _ ChatCommandPluginManagerT = (*ChatCommandPluginManager)(nil)

// NewChatCommandPluginManager creates a new ChatCommandPluginManager
func NewChatCommandPluginManager(
	ircClient client.TwitchClientT, repo repository.SingleBotRepositoryT) ChatCommandPluginManagerT {
	manager := ChatCommandPluginManager{repo: repo, pluginMap: make(map[string]chatplugins.ChatCommandPluginT)}

	manager.RegisterPlugin(commandplugins.NewAddCommandPlugin(ircClient, repo))
	manager.RegisterPlugin(commandplugins.NewDeleteCommandPlugin(ircClient, repo))
	manager.RegisterPlugin(commandplugins.NewEditCommandPlugin(ircClient, repo))
	manager.RegisterPlugin(commandplugins.NewListCommandsPlugin(ircClient, repo))
	manager.RegisterPlugin(commandplugins.NewCommandResponsePlugin(ircClient, repo))
	manager.RegisterPlugin(games.NewNumberGuesserPlugin(ircClient, repo))
	manager.RegisterPlugin(selfban.NewBanOneselfPlugin(ircClient))
	manager.RegisterPlugin(games.NewDicePlugin(ircClient))

	return &manager
}

// RegisterPlugin registers a new plugin to the manager.
func (manager *ChatCommandPluginManager) RegisterPlugin(plugin chatplugins.ChatCommandPluginT) {
	manager.pluginMap[plugin.GetPluginType()] = plugin
}

// ProcessChat parses command form chat message and calls the right plugin if necessary.
func (manager *ChatCommandPluginManager) ProcessChat(
	channel string, sender *twitch_irc.User, message *twitch_irc.PrivateMessage) {

	// Parse command name and get command, if any
	commandName := getCommandName(message.Message)
	command := manager.repo.GetCommandByChannelAndName(channel, commandName)
	if command == nil { // Chat is not a bot command.
		return
	}

	plugin, exists := manager.pluginMap[command.PluginType]
	if exists {
		plugin.ReactToChat(command, channel, sender, message)
	} else {
		// Likely synchronization issue
		log.Println("No plugin exists for plugin type " + command.PluginType + " and command " + commandName)
	}
}

func getCommandName(text string) string {
	// strings.Fields deals with heading/trailing/non-space whitespaces.
	fields := strings.Fields(text)
	if len(fields) == 0 { // This is unlikely, just checking for malicious input
		return ""
	}
	return fields[0]
}
